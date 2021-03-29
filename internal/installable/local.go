package installable

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const localRoot = "_local"

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

// GetInstallPath returns the installation path for provided local Installlable
// by joining `<installDirName>/_local/<relative-path-to-root-from-cwd>`.
// TODO decide if to accept absolute paths or only relative paths without access to parent directories
// FIXME this might be broken with certain relative paths
func (l Local) GetInstallPath() (string, error) {
	if err := l.Validate(); err != nil {
		return "", err
	}

	// calculate the relative path from abs current dir to abs l.Root
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf(`getting current working directory: %w`, err)
	}
	absRoot, err := filepath.Abs(l.Root)
	if err != nil {
		return "", fmt.Errorf(`getting absolute path to %s for local installable %+v: %w`, l.Root, l, err)
	}
	rel, err := filepath.Rel(cwd, absRoot)
	if err != nil {
		return "", fmt.Errorf(`getting relative path to %s: %w`, l.Root, err)
	}
	// if rel is a file, install path is the parent
	if info, err := os.Stat(rel); err != nil {
		return "", fmt.Errorf(`statting %s for installable %+v: %w`, rel, l, err)
	} else if !info.IsDir() {
		rel = filepath.Dir(rel)
	}

	return filepath.Join(installDirName, localRoot, rel), nil
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
