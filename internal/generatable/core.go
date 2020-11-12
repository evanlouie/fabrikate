package generatable

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	generateDirName        = "_generated"
	componentPathSeparator = "_"
)

type Generatable interface {
	// Generate the generatable to the yaml file pointed to by GetGeneratePath.
	// Returns the number of bytes written or an error.
	Generate() (int, error)
	// GetGeneratePath will convert the generatables ComponentPath into an
	// absolute path to a yaml file.
	GetGeneratePath() (string, error)
	// Validate the generatable to ensure that all fields are correct.
	// Should be as stateless as possible -- can be called multiple times
	// throughout application life (e.g. before both GetGeneratePath and Generate)
	Validate() error
}

// Pathable encapsulates the common data and methods needed for GetGeneratePath.
// All particpating struct for Generatable should inherit from this.
type Pathable struct {
	ComponentPath []string // list representation of parent->child component structure; used to generate the generation path
}

// GetGeneratePath converts the pathables ComponentPath to os-specifc filepath
// starting at `$(pwd)/_generated`
func (n Pathable) GetGeneratePath() (string, error) {
	if len(n.ComponentPath) == 0 {
		return "", fmt.Errorf(`ComponentPath must be of length >0, given %+v`, n.ComponentPath)
	}
	componentName := strings.Join(n.ComponentPath, componentPathSeparator)
	return filepath.Join(generateDirName, componentName) + ".yaml", nil
}

// cleanup removes generated components by deleting the yaml file pointed to
// GetGeneratePath
func cleanup(g Generatable) error {
	if g == nil {
		return fmt.Errorf(`nil generatable passed to cleanup`)
	}

	p, err := g.GetGeneratePath()
	if err != nil {
		return fmt.Errorf(`getting generation path for generatable %v: %w`, g, err)
	}
	if err := os.Remove(p); err != nil {
		return fmt.Errorf(`cleaning up generated generatable %v at %s: %w`, g, p, err)
	}

	return nil
}
