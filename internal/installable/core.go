package installable

import (
	"fmt"
	"os"
)

const (
	installDirName = "_components" // where all components will be installed to
)

// Installable encapsulates the functionality installing/pulling component
// resources locally.
type Installable interface {
	// Install pulls the resource specified by the Installable into the path
	// returned from GetInstallPath.
	Install() error
	// GetInstallPath returns the the installation path for the Installable
	// resource.
	// The start of the the returned path will be `$(pwd)/_components`.
	// The Installable may or may-not be Install(ed) when calling this
	// function.
	GetInstallPath() (string, error)
	// Validate the data in given Installable depending on the implementing
	// struct. This function should can be expected to be run by either external
	// from other packages and internally by other implementing functions of the
	// Installable (e.g GetInstallPath).
	Validate() error
	// Clean the installed component by removing the its fetched content. This is
	// a noop for local components.
	Clean() error
}

// clean installed components by deleting the path returned from
// GetInstallPath.
func clean(i Installable) error {
	if i == nil {
		return fmt.Errorf(`nil installable passed to cleanup`)
	}

	installPath, err := i.GetInstallPath()
	if err != nil {
		return fmt.Errorf(`getting install path for installable %v: %w`, i, err)
	}
	if err := os.RemoveAll(installPath); err != nil {
		return fmt.Errorf(`cleaning installation for installable %+v at %s: %w`, i, installPath, err)
	}

	return nil
}
