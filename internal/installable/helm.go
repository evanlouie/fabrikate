package installable

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/microsoft/fabrikate/internal/helm"
	"github.com/microsoft/fabrikate/internal/url"
)

type Helm struct {
	URL     string
	Chart   string
	Version string
}

func (h Helm) Install() error {
	if err := h.Validate(); err != nil {
		return err
	}

	// Pull to a temporary directory
	tmpHelmDir, err := os.MkdirTemp("", "fabrikate")
	defer os.RemoveAll(tmpHelmDir)
	if err != nil {
		return fmt.Errorf(`creating temporary directory to "helm pull" into: %w`, err)
	}
	if err := helm.Pull(h.URL, h.Chart, h.Version, tmpHelmDir); err != nil {
		return fmt.Errorf(`installing helm component %+v into %s: %w`, h, tmpHelmDir, err)
	}

	componentPath, err := h.GetInstallPath()
	if err != nil {
		return fmt.Errorf(`getting helm installation path for %+v: %w`, h, err)
	}
	if err := os.RemoveAll(componentPath); err != nil {
		return fmt.Errorf(`removing existing helm installation at %s: %w`, componentPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(componentPath), os.ModePerm); err != nil {
		return fmt.Errorf(`creating helm installation directory %s: %w`, componentPath, err)
	}

	// Move the extracted chart from tmp to the _component dir
	extractedChartPath := filepath.Join(tmpHelmDir, h.Chart)
	if err := os.Rename(extractedChartPath, componentPath); err != nil {
		return fmt.Errorf(`moving extracted helm chart from %s to %s: %w`, extractedChartPath, componentPath, err)
	}

	return nil
}

func (h Helm) GetInstallPath() (string, error) {
	if err := h.Validate(); err != nil {
		return "", fmt.Errorf(`validing helm installable %+v: %sw`, h, err)
	}
	urlPath, err := url.ToPath(h.URL)
	if err != nil {
		return "", fmt.Errorf(`getting install path for helm chart %+v: %w`, h, err)
	}
	var version string
	if h.Version != "" {
		version = h.Version
	} else {
		version = "latest"
	}

	componentPath := filepath.Join(installDirName, urlPath, h.Chart, version)
	return componentPath, nil
}

func (h Helm) Validate() error {
	switch {
	case h.URL == "":
		return fmt.Errorf(`URL must be non-zero length: %+v`, h)
	case h.Chart == "":
		return fmt.Errorf(`chart must be non-zero length: %+v`, h)
	}

	return nil
}

func (h Helm) Clean() error {
	return clean(h)
}
