package helm

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	fabYaml "github.com/microsoft/fabrikate/internal/yaml"
)

// TemplateOptions encapsulate the options for `helm template`.
// helm template \
//   [ --repo <Repo> ] \
//   [ --version <Version> ] \
//   [ --namespace <Namespace> --create-namespace ] \
//   [ --values <Values[0]> --values <Value[1]> ... ] \
//   [ --set <Set[0]> --set <Set[1]> ... ] \
//   [Release] <Chart>
type TemplateOptions struct {
	Release   string   // [NAME]
	Chart     string   // [CHART]
	Repo      string   // --repo
	Version   string   // --version
	Namespace string   // --namespace flag. implies --create-namespace
	Values    []string // "--value" flags. e.g.: ["foo/bar.yaml", "/etc/my/values.yaml"] == "--values foo/bar.yaml -- values /et/my/values.yaml"
	Set       []string // "--set" flags. e.g: ["foo=bar", "baz=123"] == "--set foo=bar --set baz=123"
}

// TemplateWithCRDs will `helm template` the target chart as well as ensure
// that any YAML files in the the charts "crds" directory are prepended to
// the returned YAML string -- which are not templated via "helm template" in
// helm 3.
//
// Starting with Helm 3, the "crds" directory of a chart holds a special meaning
// and holds CRD YAMLs which are not templated -- thus not outputted from
// `helm template` -- but installed to the cluster via `helm install`. This
// function is useful to get a complete YAML output for the entire chart.
func TemplateWithCRDs(opts TemplateOptions) ([]map[string]interface{}, error) {
	// interpertet the chart path based on if a repo-url was provided
	var chartPath, crdPath string
	if opts.Repo != "" {
		tmpDir, err := os.MkdirTemp("", "fabrikate")
		if err != nil {
			return nil, fmt.Errorf(`creating temporary directory to pull helm chart %s@%s from %s: %w`, opts.Chart, opts.Version, opts.Repo, err)
		}
		defer os.RemoveAll(tmpDir)
		if err := Pull(opts.Repo, opts.Chart, opts.Version, tmpDir); err != nil {
			return nil, fmt.Errorf(`pulling helm chart %s@%s from %s: %w`, opts.Chart, opts.Version, opts.Repo, err)
		}
		chartPath = filepath.Join(tmpDir, opts.Chart)
	} else {
		chartPath = opts.Chart
	}
	crdPath = filepath.Join(chartPath, "crds")

	// walk the "crds" dir to collect all the yaml strings
	var crds []string // list of crd yaml <strings>
	if info, err := os.Stat(crdPath); err == nil {
		if info.IsDir() {
			err := filepath.Walk(crdPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf(`walking CRD path %s: %w`, path, err)
				}
				extension := strings.ToLower(filepath.Ext(info.Name()))
				// track all yaml files
				if !info.IsDir() && extension == ".yaml" {
					crd, err := os.ReadFile(path)
					if err != nil {
						return fmt.Errorf("reading CRD file %s: %w", path, err)
					}
					crds = append(crds, string(crd))
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf(`walking CRD path %s: %w`, crdPath, err)
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf(`reading helm chart CRD directory %s: %w`, crdPath, err)
	}

	// run `helm template` to get the contents of the pulled chart
	templateOpts := opts           // inherit all the initial settings
	templateOpts.Repo = ""         // zero out so it wont attempt to lookup the repo
	templateOpts.Chart = chartPath // manually set the path of the chart to the downloaded chart
	templateMaps, err := Template(templateOpts)
	if err != nil {
		return nil, fmt.Errorf(`templating helm chart at %s: %w`, templateOpts.Chart, err)
	}

	// decode CRDs to maps
	crdMaps, err := fabYaml.DecodeMaps([]byte(strings.Join(crds, "\n---\n")))
	if err != nil {
		return nil, fmt.Errorf(`decoding CRD yaml: %w`, err)
	}

	// merge to unified map list
	var finalYamlMaps []map[string]interface{}
	for _, crdMap := range crdMaps {
		if crdMap != nil {
			finalYamlMaps = append(finalYamlMaps, crdMap)
		}
	}
	for _, templateMap := range templateMaps {
		if templateMap != nil {
			finalYamlMaps = append(finalYamlMaps, templateMap)
		}
	}

	return finalYamlMaps, nil
}

// Template runs `helm template` on the chart specified by opts.
// Returns the map[string]interface{} decoded output of stdout for the
// `helm template` call.
// Will return an error occurs when:
//
//  - running `helm template` returns an error
//  - command outputs ANYTHING to stderr
//  - stdout of `helm template` command cannot be decoded to a map[string]interface{}
//
// NOTE in Helm 3, CRDs in the "crds" directory of the chart are not outputted
// from `helm template` but are installed via `helm install`
func Template(opts TemplateOptions) ([]map[string]interface{}, error) {
	templateArgs := []string{"template"}
	if opts.Repo != "" {
		// if an existing helm repo exists on the helm client, use that for templating
		existingRepo, err := FindRepoNameByURL(opts.Repo)
		if err != nil {
			return nil, fmt.Errorf(`searching existing helm repositories for %s: %w`, opts.Repo, err)
		}
		if existingRepo != "" {
			opts.Chart = path.Join(existingRepo, opts.Chart)
		} else {
			// if an existing repo is not found, use the --repo option to pull from network
			templateArgs = append(templateArgs, "--repo", opts.Repo)
		}
	}
	// set namespace if provided
	if opts.Namespace != "" {
		templateArgs = append(templateArgs, "--create-namespace", "--namespace", opts.Namespace)
	}
	// set all --set options
	for _, set := range opts.Set {
		templateArgs = append(templateArgs, "--set", set)
	}
	// set all --values options
	for _, yamlPath := range opts.Values {
		templateArgs = append(templateArgs, "--values", yamlPath)
	}

	// a helm release [NAME] is specified as an optional leading parameter to the [CHART]
	if opts.Release != "" {
		templateArgs = append(templateArgs, opts.Release)
	}
	templateArgs = append(templateArgs, opts.Chart)

	templateCmd := exec.Command("helm", templateArgs...)
	var stdout, stderr bytes.Buffer
	templateCmd.Stdout = &stdout
	templateCmd.Stderr = &stderr

	if err := templateCmd.Run(); err != nil {
		return nil, fmt.Errorf(`running "%s": %s: %w`, templateCmd, stderr.String(), err)
	}
	if stderr.Len() != 0 {
		return nil, fmt.Errorf(`"%s" exited with output to stderr: %s`, templateCmd, stderr.String())
	}
	maps, err := fabYaml.DecodeMaps(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf(`parsing output of "helm template": %s: %w`, stderr.String(), err)
	}

	return maps, nil
}

func injectNamespace(manifest map[string]interface{}, namespace string) (map[string]interface{}, error) {
	if manifest == nil {
		return nil, nil
	}
	// inject the metadata map if it is not present
	if _, ok := manifest["metadata"]; !ok {
		manifest["metadata"] = map[string]interface{}{}
	}

	metadata, ok := manifest["metadata"].(map[string]interface{})
	switch {
	case !ok:
		return nil, fmt.Errorf(`asserting metadata of yaml manifest is map[string]interface{}: %+v`, manifest)
	case metadata["namespace"] != nil:
		return nil, fmt.Errorf(`existing namespace found in yaml: %+v`, manifest)
	default:
		metadata["namespace"] = namespace
	}

	return manifest, nil
}

func createNamespace(name string) map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": name,
		},
	}
}
