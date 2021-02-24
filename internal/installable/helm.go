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
		return err
	}
	if err := os.RemoveAll(componentPath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(componentPath), os.ModePerm); err != nil {
		return err
	}

	// Move the extracted chart from tmp to the _component dir
	extractedChartPath := filepath.Join(tmpHelmDir, h.Chart)
	if err := os.Rename(extractedChartPath, componentPath); err != nil {
		return err
	}

	return nil
}

func (h Helm) GetInstallPath() (string, error) {
	if err := h.Validate(); err != nil {
		return "", err
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
	if h.URL == "" {
		return fmt.Errorf(`URL must be non-zero length: %+v`, h)
	}
	if h.Chart == "" {
		return fmt.Errorf(`Chart must be non-zero length: %+v`, h)
	}

	return nil
}
