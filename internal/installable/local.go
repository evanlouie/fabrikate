package installable

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// track all files
	var files []string

	// if targeting a directory, walk it and add all files to tracking list
	// otherwise, just add the Root to tracking list
	if targetInfo.IsDir() {
		err = filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf(`walking local installable path %s: %w`, path, err)
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
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

	// copy all files to the installPath
	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			return fmt.Errorf(`opening file %s for local installable %+v: %w`, file, l, err)
		}
		defer in.Close()
		copyPath := filepath.Join(installPath, filepath.Base(file))
		if err := os.MkdirAll(filepath.Dir(copyPath), os.ModePerm); err != nil {
			return fmt.Errorf(`creating parent directory for file %s: %w`, copyPath, err)
		}
		out, err := os.Create(copyPath)
		if err != nil {
			return fmt.Errorf(`creating file %s for copying of contents for local installable %+v: %w`, copyPath, l, err)
		}
		defer out.Close()
		if _, err := io.Copy(out, in); err != nil {
			return fmt.Errorf(`copying contents from %s to %s for local installabled %+v: %w`, file, copyPath, l, err)
		}
	}

	return nil
}

// GetInstallPath returns the installation path for provided local Installlable
// by joining
// `<installDirName>/<absolute-path-to-root-with-seperators-replaced-with-$>`.
//
// TODO decide if to accept absolute paths or only relative paths without access to parent directories
// FIXME this might be broken with certain relative paths
func (l Local) GetInstallPath() (string, error) {
	if err := l.Validate(); err != nil {
		return "", err
	}

	absRoot, err := filepath.Abs(l.Root)
	if err != nil {
		return "", fmt.Errorf(`getting absolute path to %s for local installable %+v: %w`, l.Root, l, err)
	}

	// if the target is a file, the installable directory will be the parent
	if info, err := os.Stat(absRoot); err != nil {
		return "", fmt.Errorf(`statting %s for installable %+v: %w`, absRoot, l, err)
	} else if !info.IsDir() {
		absRoot = filepath.Dir(absRoot)
	}
	flattendAbsPath := strings.ReplaceAll(absRoot, string(filepath.Separator), "$")

	return filepath.Join(installDirName, flattendAbsPath), nil
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
		return fmt.Errorf(`local installable root must be within or a child of current working directory %s: resolved %s`, cwd, absRoot)
	}

	return nil
}
