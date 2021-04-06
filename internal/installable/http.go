package installable

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/microsoft/fabrikate/internal/url"
	"github.com/microsoft/fabrikate/internal/yaml"
)

type HTTP struct {
	URL string
}

func (h HTTP) Install() error {
	if err := h.Validate(); err != nil {
		return err
	}

	// deleting if it already exists
	componentPath, err := h.GetInstallPath()
	if err != nil {
		return fmt.Errorf(`getting install path for http component %+v: %w`, h, err)
	}
	if err := os.RemoveAll(componentPath); err != nil {
		return fmt.Errorf(`cleaning previous http component installation at "%s": %w`, componentPath, err)
	}

	// download the resource
	resp, err := http.Get(h.URL)
	if err != nil {
		return fmt.Errorf(`fetching URL %s for HTTP installable: %w`, h.URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(`reading return bytes for HTTP installable: %w`, err)
	}

	// validate that the contents is yaml by parsing it
	if _, err := yaml.Decode(body); err != nil {
		return fmt.Errorf(`response for HTTP installable %s is not valid YAML: %w`, h.URL, err)
	}

	parentDir := filepath.Dir(componentPath)
	if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return fmt.Errorf(`creating install directory %s for HTTP installable %s: %w`, parentDir, h.URL, err)
	}
	if err := os.WriteFile(componentPath, body, os.ModePerm); err != nil {
		return fmt.Errorf(`writing fetched yaml document for HTTP installable %s to %s: %w`, h.URL, componentPath, err)
	}

	return nil
}

func (h HTTP) GetInstallPath() (string, error) {
	if err := h.Validate(); err != nil {
		return "", err
	}

	urlPath, err := url.ToPath(h.URL)
	if err != nil {
		return "", fmt.Errorf(`getting install path for helm chart %+v: %w`, h, err)
	}

	return filepath.Join(installDirName, urlPath), nil
}

func (h HTTP) Validate() error {
	httpProtocolRgx := regexp.MustCompile(`(?i)^https?://.+$`)
	if !httpProtocolRgx.MatchString(h.URL) {
		return fmt.Errorf(`HTTP installable URL %s does not match validating regex %s`, h.URL, httpProtocolRgx)
	}

	return nil
}
