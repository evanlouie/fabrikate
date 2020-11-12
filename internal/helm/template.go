package helm

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// TemplateOptions encapsulate the options for `helm template`
type TemplateOptions struct {
	Release   string
	RepoURL   string
	Chart     string
	Version   string
	Namespace string
	Values    []string
}

// Template is a command for `helm template`
func Template(opts TemplateOptions) (string, error) {
	templateArgs := []string{"template"}
	if opts.Release != "" {
		templateArgs = append(templateArgs, "--release-name", opts.Release)
	}
	if opts.RepoURL != "" {
		templateArgs = append(templateArgs, "--repo", opts.RepoURL)
	}
	if opts.Namespace != "" {
		templateArgs = append(templateArgs, "--create-namespace", "--namespace", opts.Namespace)
	}
	for _, yamlPath := range opts.Values {
		templateArgs = append(templateArgs, "--values", yamlPath)
	}
	templateArgs = append(templateArgs, opts.Chart)

	templateCmd := exec.Command("helm", templateArgs...)
	var stdout, stderr bytes.Buffer
	templateCmd.Stdout = &stdout
	templateCmd.Stderr = &stderr

	if err := templateCmd.Run(); err != nil {
		return "", fmt.Errorf(`running "%s": %v: %v`, templateCmd, err, stderr.String())
	}

	return stdout.String(), nil
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
	if !ok {
		return nil, fmt.Errorf(`reflecting metadata of yaml manifest: %+v`, manifest)
	}
	if metadata["namespace"] != nil {
		return nil, fmt.Errorf(`existing namespace found in yaml: %+v`, manifest)
	}
	metadata["namespace"] = namespace

	return manifest, nil
}

func injectNamespaceBack(unifiedManifest string, namespace string) (string, error) {
	// split the unified manifest string by "---"
	dividerRgx := regexp.MustCompile(`^---$`)
	manifestStrings := dividerRgx.Split(unifiedManifest, -1)

	// parse and inject the namespace into the parsed map
	var injectedManifests []string
	for _, entry := range manifestStrings {
		var m map[interface{}]interface{}
		if err := yaml.Unmarshal([]byte(entry), &m); err != nil {
			return "", fmt.Errorf(`unmarshalling YAML string %s: %w`, entry, err)
		}
		if m["metadata"] != nil {
			metadata, ok := m["metadata"].(map[string]interface{})
			if !ok {
				return "", fmt.Errorf(`reflecting metadata of yaml manifest: %+v`, m)
			}
			if metadata["namespace"] == nil {
				metadata["namespace"] = namespace
			}
		}
		asBytes, err := yaml.Marshal(m)
		if err != nil {
			return "", fmt.Errorf(`marshalling namespace injected YAML %+v: %w`, m, err)
		}
		injectedManifests = append(injectedManifests, string(asBytes))
	}

	// re-join the strings with "---"
	withNS := strings.TrimSpace(strings.Join(injectedManifests, "\n---\n"))

	return strings.TrimSpace(withNS), nil
}

// cleanManifest parses either a yaml document (or list of documents delimitted
// by "---") and removes entries that are not of type map[string]interface{}.
//
// TODO find out if this is needed in helm 3
func cleanManifest(manifest string) (string, error) {
	// split based on yaml divider
	manifests := strings.Split(manifest, "\n---")

	// remove all invalid yaml
	var cleaned []string
	for _, entry := range manifests {
		var m map[string]interface{}
		// if it doesn't unmarshal properly, do not add
		if err := yaml.Unmarshal([]byte(entry), &m); err != nil {
			continue
		}
		// only append documents with a non-empty body
		if len(strings.TrimSpace(entry)) > 0 {
			cleaned = append(cleaned, entry)
		}
	}

	// re-join based on yaml divider
	joined := strings.TrimSpace(strings.Join(cleaned, "\n---"))
	return joined, nil
}
