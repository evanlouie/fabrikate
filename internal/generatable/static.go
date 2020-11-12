package generatable

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Static struct {
	Pathable
	ManifestPath string // path to static manifests
}

func (s Static) Validate() error {
	if _, err := os.Stat(s.ManifestPath); os.IsNotExist(err) {
		return fmt.Errorf(`ManifestPath for static generator does not exist: %+v`, s)
	}
	if len(s.ComponentPath) == 0 {
		return fmt.Errorf(`ComponentPath for static generator must be of length >0: %+v`, s)
	}

	return nil
}

func (s Static) Generate() (int, error) {
	if err := s.Validate(); err != nil {
		return 0, fmt.Errorf(`invalid static generator %+v: %w`, s, err)
	}

	// Load all manifest strings in ManifestPath
	var manifests []string
	yamlExtRgx := regexp.MustCompile(`(?i)\.ya?ml$`) // helper regex to see if a file is a yaml file
	err := filepath.Walk(s.ManifestPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf(`walking static manifest path "%s": %w`, path, err)
		}
		// if it is a yaml file, add it to the manifest list
		if !info.IsDir() && yamlExtRgx.MatchString(filepath.Ext(path)) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf(`reading YAML file at "%s": %w`, path, err)
			}
			manifests = append(manifests, string(b))
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf(`reading YAML files in "%s": %w`, s.ManifestPath, err)
	}

	// delete existing generation path
	generatePath, err := s.GetGeneratePath()
	if err != nil {
		return 0, fmt.Errorf(`getting generation path: %w`, err)
	}
	if err := os.RemoveAll(generatePath); err != nil {
		return 0, fmt.Errorf(`cleaning existing static component generation at %s: %w`, generatePath, err)
	}

	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(generatePath), os.ModePerm); err != nil {
		return 0, fmt.Errorf(`creating parent directory %s for static component generation: %w`, generatePath, err)
	}

	// Write manifests to generation path
	unifiedManifest := strings.Join(manifests, "\n---\n")
	asBytes := []byte(unifiedManifest)
	if err := ioutil.WriteFile(generatePath, asBytes, os.ModePerm); err != nil {
		return 0, fmt.Errorf(`writing generated static component to %s: %w`, generatePath, err)
	}

	return len(asBytes), nil
}
