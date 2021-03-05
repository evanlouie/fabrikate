package installable

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	Root string
}

// Install will attempt to copy all folders/files that the local installable
// coordinate points to the
func (l Local) Install() error {
	if err := l.Validate(); err != nil {
		return err
	}

	abs, err := filepath.Abs(l.Root)
	if err != nil {
		return fmt.Errorf(`computing absolute path for local installable %+v: %w`, l, err)
	}
	targetInfo, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf(`stating Root for local component %+v: %w`, l, err)
	}

	// track all folders and files
	var folders, files []string

	// if targeting a directory, walk it.
	// otherwise, just add the Root to files and its parent directory to folders
	if targetInfo.IsDir() {
		err = filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf(`walking local installable path %s: %w`, path, err)
			}
			if info.IsDir() {
				folders = append(folders, path)
			} else {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		folders = append(folders, filepath.Dir(abs))
		files = append(files, abs)
	}

	// clean existing installation
	installPath, err := l.GetInstallPath()
	if err != nil {
		return fmt.Errorf(`getting install path for local installable %+v: %w`, l, err)
	}
	if err := os.RemoveAll(installPath); err != nil {
		return fmt.Errorf(`cleaning existing local component installation: %w`, err)
	}

	// copy all folders and files to the installPath; folders go first
	for _, dir := range folders {
		folderPath := filepath.Join(installPath, dir)
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return fmt.Errorf(`copying directory %s for local installable %+v: %w`, dir, l, err)
		}
	}
	for _, file := range files {
		originalPath, err := filepath.Abs(file)
		if err != nil {
			return fmt.Errorf(`computing absolute path for copied resource %s from installable %+v: %w`, originalPath, l, err)
		}
		in, err := os.Open(originalPath)
		if err != nil {
			return fmt.Errorf(`opening file %s for local installable %+v: %w`, originalPath, l, err)
		}
		defer in.Close()
		copyPath := filepath.Join(installPath, file)
		out, err := os.Create(copyPath)
		if err != nil {
			return fmt.Errorf(`creating file %s for copying of contents for local installable %+v: %w`, copyPath, l, err)
		}
		defer out.Close()
		if _, err := io.Copy(out, in); err != nil {
			return fmt.Errorf(`copying contents from %s to %s for local installabled %+v: %w`, originalPath, copyPath, l, err)
		}
	}

	return nil
}

// GetInstallPath returns the absolute path of the local Installables Root and
// prepends it with the component installation directory.
func (l Local) GetInstallPath() (string, error) {
	if err := l.Validate(); err != nil {
		return "", err
	}

	var resolvedAbsolute string
	if filepath.IsAbs(l.Root) {
		resolvedAbsolute = l.Root
	} else {
		abs, err := filepath.Abs(l.Root)
		if err != nil {
			return "", fmt.Errorf(`computing absolute path for %s: %w`, l.Root, err)
		}
		resolvedAbsolute = abs
	}

	info, err := os.Stat(resolvedAbsolute)
	if err != nil {
		return "", fmt.Errorf(`stating %s for local component %+v: %w`, resolvedAbsolute, l, err)
	}
	if !info.IsDir() {
		resolvedAbsolute = filepath.Dir(resolvedAbsolute)
	}

	return filepath.Join(installDirName, resolvedAbsolute), nil
}

func (l Local) Validate() error {
	switch {
	case l.Root == "":
		return fmt.Errorf(`Root must be non-zero length: %+v`, l)
	}

	// ensure the root exists
	if _, err := os.Stat(l.Root); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(`Root does not exist in filesystem for local installable %+v: %w`, l, err)
		}
		return fmt.Errorf(`unable to stat Root for %+v: %w`, l, err)
	}

	return nil
}
