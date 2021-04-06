package installable

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	Root string
}

// Install is a noop for local installables.
func (l Local) Install() error {
	if err := l.Validate(); err != nil {
		return err
	}

	return nil
}

// GetInstallPath returns the path where the local installable is located on
// disk.
func (l Local) GetInstallPath() (string, error) {
	if err := l.Validate(); err != nil {
		return "", err
	}

	return l.Root, nil
}

func (l Local) Validate() error {
	if l.Root == "" {
		return fmt.Errorf(`local installable root must be non-zero length: %+v`, l)
	}

	// ensure the root exists
	if _, err := os.Stat(l.Root); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(`local installable root does not exist in filesystem for local installable %+v: %w`, l, err)
		}
		return fmt.Errorf(`unable to stat Root for %+v: %w`, l, err)
	}

	// root must be in the component file tree -- cwd must be the prefix of the abs l.root
	absRoot, err := filepath.Abs(l.Root)
	if err != nil {
		return fmt.Errorf(`getting absolute path of %s: %w`, l.Root, err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(`getting current working directory: %w`, err)
	}
	if !strings.HasPrefix(absRoot, cwd) {
		return fmt.Errorf(`local installable root must be subdirectory of the current working directory %s: resolved %s`, cwd, absRoot)
	}

	return nil
}

// Clean is a noop for local components.
func (l Local) Clean() error {
	return nil
}
