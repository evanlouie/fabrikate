package installable

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/microsoft/fabrikate/internal/url"
)

type Git struct {
	URL    string
	SHA    string
	Branch string
}

func (g Git) Install() error {
	if err := g.Validate(); err != nil {
		return err
	}

	// deleting if it already exists
	componentPath, err := g.GetInstallPath()
	if err != nil {
		return fmt.Errorf(`getting install path for git component %+v: %w`, g, err)
	}
	if err := os.RemoveAll(componentPath); err != nil {
		return fmt.Errorf(`cleaning previous git component installation at "%s": %w`, componentPath, err)
	}

	if err := g.clone(componentPath); err != nil {
		return fmt.Errorf(`cloning git component from "%s" into "%s": %w`, g.URL, componentPath, err)
	}

	return nil
}

func (g Git) GetInstallPath() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}
	urlPath, err := url.ToPath(g.URL)
	if err != nil {
		return "", fmt.Errorf(`generating installation path for git component %+v: %w`, g, err)
	}

	var version string
	if g.SHA != "" {
		version = g.SHA
	} else if g.Branch != "" {
		version = g.Branch
	} else {
		version = "latest"
	}

	componentPath := filepath.Join(installDirName, urlPath, version)
	return componentPath, nil
}

func (g Git) Validate() error {
	switch {
	case g.URL == "":
		return fmt.Errorf(`URL must be non-zero length`)
	case g.SHA != "" && g.Branch != "":
		return fmt.Errorf(`Only one of SHA or Branch can be provided, "%v" and "%v" provided respectively`, g.SHA, g.Branch)
	}

	return nil
}

//------------------------------------------------------------------------------
// Git Helpers
//------------------------------------------------------------------------------

var cloneCoordinator = struct {
	sync.Mutex                          // lock to ensure only one write has access to locks at a time
	nodes      map[string]*sync.RWMutex // key == filepath; value == lock denoting if the filepath has been clone or is cloning
}{
	nodes: map[string]*sync.RWMutex{},
}

// clone performs a `git clone <g.URL> <dir>`
func (g Git) clone(dir string) error {
	nodes := cloneCoordinator.nodes
	cloneCoordinator.Lock() // establish a lock so we can safely read from the map of locks

	// If one exists, another thread is cloning it; just need to wait for it to
	// become free; establish a lock and immediately release
	if node, exists := nodes[dir]; exists {
		node.RLock()
		defer node.RUnlock()
		cloneCoordinator.Unlock()
		return nil
	}

	// It is possible that another channel attempted to create the same mutex and
	// established a lock before this one. Do a final check to see if a lock exists
	if _, exists := nodes[dir]; exists {
		return nil
	}

	// create a mutex for the path
	nodes[dir] = &sync.RWMutex{} // add a rwlock

	node, exists := nodes[dir]
	if !exists {
		return fmt.Errorf(`error creating mutex lock for cloning repo "%v" to dir "%v"`, g.URL, dir)
	}

	// write lock the path to block others from cloning the same path
	node.Lock() // establish a write lock so the other readers are blocked
	defer node.Unlock()
	cloneCoordinator.Unlock()

	// Prep the clone args
	cloneOpts := git.CloneOptions{
		URL: g.URL,
		// Progress: os.Stdout, //  TODO encapsulate in a feature flag
	}
	// add the branch to clone options if present
	if g.Branch != "" {
		cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(g.Branch)
		cloneOpts.Depth = 1
	}

	// clone the repo
	r, err := git.PlainClone(dir, false, &cloneOpts)
	if err != nil {
		return fmt.Errorf(`cloning git repository "%s" into "%s": %w`, g.URL, dir, err)
	}

	// checkout the SHA if provided
	if g.SHA != "" {
		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf(`getting worktree for git repository %s: %w`, dir, err)
		}
		if err := w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(g.SHA),
		}); err != nil {
			return fmt.Errorf(`checking out SHA "%s" in git repo at "%s": %w`, g.SHA, dir, err)
		}
	}

	// ensure the target branch or SHA is checked out
	head, err := r.Head()
	if err != nil {
		return fmt.Errorf(`getting HEAD of repo at %s: %w`, dir, err)
	}
	switch {
	case g.SHA != "":
		if head.Hash().String() != g.SHA {
			return fmt.Errorf(`repo at %s not checked out to target SHA %s: is at %s`, dir, g.SHA, head.Hash())
		}
	case g.Branch != "":
		if !strings.HasSuffix(string(head.Name()), g.Branch) {
			return fmt.Errorf(`repo at %s not checked out to target branch %s: is at %s`, dir, g.Branch, head.Name())
		}
	}

	return nil
}
