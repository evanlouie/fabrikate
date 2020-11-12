package generatable

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/microsoft/fabrikate/internal/helm"
	"github.com/timfpark/yaml"
)

type Helm struct {
	Pathable
	ChartPath string // path to directory containing Chart.yaml. Not to actual chart.yaml file
	Values    map[string]interface{}
}

func (h Helm) Validate() error {
	return nil
}

func (h Helm) Generate() (int, error) {
	// write values.yaml to a temporary file
	valuesFile, err := ioutil.TempFile("", "fabrikate")
	if err != nil {
		return 0, fmt.Errorf(`creating temporary helm values files: %w`, err)
	}
	defer os.Remove(valuesFile.Name())
	valueBytes, err := yaml.Marshal(h.Values)
	if err != nil {
		return 0, fmt.Errorf(`marshalling helm values to temporary file: %w`, err)
	}
	if _, err := valuesFile.Write(valueBytes); err != nil {
		return 0, fmt.Errorf(`writing temporary helm values file: %w`, err)
	}
	if err := valuesFile.Close(); err != nil {
		return 0, fmt.Errorf(`writing temporary helm values file: %w`, err)
	}

	generatePath, err := h.GetGeneratePath()
	if err != nil {
		return 0, fmt.Errorf(`getting generation path for helm component %+v: %w`, h, err)
	}

	// release will be set to the generate yaml file name without ext
	release := func() string {
		base := filepath.Base(generatePath)
		baseNoExt := strings.TrimSuffix(base, filepath.Ext(base))
		return baseNoExt
	}()

	// run `helm template`
	template, err := helm.Template(helm.TemplateOptions{
		Chart:   h.ChartPath,
		Values:  []string{valuesFile.Name()},
		Release: release,
	})
	if err != nil {
		return 0, fmt.Errorf(`helm template error: %w`, err)
	}

	// remove existing generation
	if err := os.Remove(generatePath); err != nil {
		return 0, fmt.Errorf(`removing previously generated helm component at "%s": %w`, generatePath, err)
	}

	// write out template
	asBytes := []byte(template)
	if err := ioutil.WriteFile(generatePath, asBytes, os.ModePerm); err != nil {
		return 0, fmt.Errorf(`writing out helm component at "%s": %w`, generatePath, err)
	}

	return len(asBytes), nil
}
