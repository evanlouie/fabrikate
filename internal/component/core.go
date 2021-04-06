package component

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microsoft/fabrikate/internal/component/config"
	"github.com/microsoft/fabrikate/internal/generatable"
	"github.com/microsoft/fabrikate/internal/installable"
	"gopkg.in/yaml.v3"
)

type hooks struct {
	BeforeInstall  []string `json:"before-install,omitempty" yaml:"before-install,omitempty"`
	AfterInstall   []string `json:"after-install,omitempty" yaml:"after-install,omitempty"`
	BeforeGenerate []string `json:"before-generate,omitempty" yaml:"before-generate,omitempty"`
	AfterGenerate  []string `json:"after-generate,omitempty" yaml:"after-generate,omitempty"`
}

type Component struct {
	Name    string
	Hooks   hooks
	Kind    string `yaml:"type,omitempty" json:"type,omitempty"` // generatable
	Method  string // installable
	Source  string // installable: git, helm, local
	Path    string // installable: helm, local
	Version string // git, helm
	Branch  string // git

	Subcomponents []Component

	config config.Config

	logicalPath  string
	physicalPath string

	manifest string
}

// Load attempts to load a `component\.(ya?ml|json)` from the provided component
// directory.
func Load(componentDirectory string, parentLogicalPath string) (c Component, err error) {
	fileRgx := regexp.MustCompile(`(?i)^component\.(ya?ml|json)$`)
	if file, err := os.Stat(componentDirectory); err != nil {
		return c, fmt.Errorf(`loading component from "%s": %w`, componentDirectory, err)
	} else if !file.IsDir() {
		return c, fmt.Errorf(`component must be a directory containing a file matching "%s", a non-directory was passed: %s`, fileRgx, componentDirectory)
	}

	// calculate possible component.<yaml|json> filepaths
	componentAbsDir, err := filepath.Abs(componentDirectory)
	if err != nil {
		return c, fmt.Errorf(`finding absolute path of component "%v": %w`, componentDirectory, err)
	}
	files, err := os.ReadDir(componentAbsDir)
	if err != nil {
		return c, fmt.Errorf(`loading files from component directory "%s": %w`, componentAbsDir, err)
	}

	// iterate over all files in the the component directory and search for a valid component.(ya?ml|json)
	// TODO refactor to just lookup the 3 files specifically instead of searching the entire directory
	var componentFile string
	for _, file := range files {
		if fileRgx.MatchString(file.Name()) {
			// make sure only one of .yaml,yml,json exists
			if componentFile == "" {
				componentFile = file.Name()
			} else {
				return c, fmt.Errorf(`only one of component definition can exist per component directory, both "%s" and "%s" found in "%s"`, componentFile, file.Name(), componentAbsDir)
			}
		}
	}
	if componentFile == "" {
		return c, fmt.Errorf(`component definition file matching %s not found in "%s"`, fileRgx, componentAbsDir)
	}

	// read the file and unmarshal into found extension type
	componentPath := filepath.Join(componentAbsDir, componentFile)
	componentBytes, err := os.ReadFile(componentPath)
	if err != nil {
		return c, fmt.Errorf(`failed to read component file "%v": %w`, componentPath, err)
	}

	ext := filepath.Ext(componentFile)
	switch {
	case regexp.MustCompile(`(?i)^\.ya?ml$`).MatchString(ext):
		if err := yaml.Unmarshal(componentBytes, &c); err != nil {
			return c, fmt.Errorf(`failed to unmarshal %s at "%s": %w`, componentFile, componentPath, err)
		}
	case strings.EqualFold(ext, ".json"):
		if err := json.Unmarshal(componentBytes, &c); err != nil {
			return c, fmt.Errorf(`failed to unmarshal %s at "%s": %w`, componentFile, componentPath, err)
		}
	default:
		return c, fmt.Errorf(`invalid component extension "%s"`, ext)
	}

	c.physicalPath = filepath.Clean(componentDirectory)
	c.logicalPath = path.Join(parentLogicalPath, c.Name)

	return c, err
}

func (c Component) ToInstallable() (installer installable.Installable, err error) {
	switch c.Method {
	case "git":
		installer = installable.Git{
			URL:    c.Source,
			Branch: c.Branch,
			SHA:    c.Version,
		}
	case "helm":
		installer = installable.Helm{
			URL:     c.Source,
			Chart:   c.Path,
			Version: c.Version,
		}
	case "local":
		installer = installable.Local{
			Root: filepath.Join(c.physicalPath, c.Source),
		}
	case "http":
		installer = installable.HTTP{
			URL: c.Source,
		}
	case "":
		// noop
	default:
		return installer, fmt.Errorf(`unsupported method "%s" in component "%+v"`, c.Method, c)
	}

	return installer, nil
}

func (c Component) ToGeneratable() (generator generatable.Generatable, err error) {
	installer, err := c.ToInstallable()
	if err != nil {
		return generator, fmt.Errorf(`converting to component %+v to Installable: %w`, c, err)
	}
	var installPath string
	if installer != nil {
		installPath, err = installer.GetInstallPath()
		if err != nil {
			return generator, fmt.Errorf(`getting install path for component %+v: %w`, c, err)
		}
	}
	// TODO THIS IS BUSTED
	pathable := generatable.Pathable{
		ComponentPath: strings.Split(c.logicalPath, "/"),
	}

	switch c.Kind {
	case "helm":
		generator = generatable.Helm{
			ChartPath: installPath,
			Pathable:  pathable,
		}
	case "static":
		generator = generatable.Static{
			ManifestPath: installPath,
			Pathable:     pathable,
		}
	case "":
		fallthrough // same as "component"
	case "component":
		// noop
	default:
		return generator, fmt.Errorf(`unsupported type "%s" in component %+v`, c.Kind, c)
	}

	return generator, nil
}

// Validate returns an error if the component is unable to be converted to an
// Installable and Generatable and have both of them be validated.
func (c Component) Validate() error {
	componentNameRgx := regexp.MustCompile(`(?i)^[^/]+$`)
	if !componentNameRgx.MatchString(c.Name) {
		return fmt.Errorf(`invalid component names %s: component name must match regex %s`, c.Name, componentNameRgx)
	}

	// installer, err := c.ToInstallable()
	// if err != nil {
	// 	return fmt.Errorf(`converting component %+v to installable: %w`, c, err)
	// }
	// if err := installer.Validate(); err != nil {
	// 	return fmt.Errorf(`validating installable %+v for component %+v: %w`, installer, c, err)
	// }
	// generator, err := c.ToGeneratable()
	// if err != nil {
	// 	return fmt.Errorf(`converting component %+v to generatable: %w`, c, err)
	// }
	// if err := generator.Validate(); err != nil {
	// 	return fmt.Errorf(`validaing generatable %+v for component %+v: %w`, generator, c, err)
	// }

	return nil
}

func echo(level int, message interface{}) {
	decorator := "-"
	switch level {
	case 0:
		decorator = ">"
	case 1:
		decorator = "\u2192" // right arrow
	case 2:
		decorator = "+"
	}
	// indent := strings.Repeat("    ", level)
	indent := strings.Repeat("\t", level)
	fmt.Printf("%v%v %v\n", indent, decorator, message)
}

func Install(startPath string) ([]string, error) {
	echo(0, fmt.Sprintf(`Starting Fabrikate installation at: "%v"`, startPath))
	c, err := Load(startPath, "")
	if err != nil {
		return nil, fmt.Errorf(`loading component from path '%s': %w`, startPath, err)
	}

	visited, err := c.install()
	if err != nil {
		return nil, err
	}

	echo(0, "Installation report:")
	type installReport struct {
		Notes []string `json:"_notes,omitempty"`
		// TODO decide whether to hold this in a map or a list of {Name: string, Path: string}; map iteration is not stable in Go
		Components map[string]string // logicalPath => physicalPath
	}

	report := installReport{
		Notes: []string{
			"This file is auto generated via `fab install`",
			"This file is consumed by `fab generate`",
			"The API for this file is unstable -- do not build tooling around it",
			"Order of components matters",
			"This files location relative to directory where `fab install` was called matters -- do not move it",
			"Do not modify unless you know what you are doing!",
		},
		Components: map[string]string{},
	}
	var installedLogicalPaths []string
	for _, c := range visited {
		// add to installation report
		if existingPath, ok := report.Components[c.logicalPath]; ok {
			return nil, fmt.Errorf(`duplicate installation reported for component with logical path %s: existing %s: new %s`, c.logicalPath, existingPath, c.physicalPath)
		}
		report.Components[c.logicalPath] = c.physicalPath

		installedLogicalPaths = append(installedLogicalPaths, c.logicalPath)
		echo(1, fmt.Sprintf("%s => %s", c.logicalPath, c.physicalPath))
	}

	installReportPath := "_install.lock.json"
	echo(0, fmt.Sprintf("Writing installation report to %s", installReportPath))
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return installedLogicalPaths, fmt.Errorf(`marshalling installation report %s: %w`, installReportPath, err)
	}
	if err := os.WriteFile(installReportPath, b, os.ModePerm); err != nil {
		return installedLogicalPaths, fmt.Errorf(`writing installation report %s: %w`, installReportPath, err)
	}

	return installedLogicalPaths, nil
}

func (c Component) install() (visited []Component, err error) {
	queue := []Component{c}

	//============================================================================
	// Recursive loop
	for {
		//--------------------------------------------------------------------------
		// base case
		if len(queue) == 0 {
			return visited, nil
		}

		//--------------------------------------------------------------------------
		// recursive case
		first, rest := queue[0], queue[1:]

		echo(1, fmt.Sprintf(`Installing component: "%s"`, first.logicalPath))
		echo(2, "Validating component")
		if err := c.Validate(); err != nil {
			return visited, err
		}
		echo(2, "Adding subcomponents to queue")
		for _, sub := range first.Subcomponents {
			// manually set the subcomponent paths.
			// - mimic the behavior of Load() for logicalPath
			// - set the physicalPath to the parent; will be overwritten during install if the component has a valid Installable
			// TODO investigate if there is a cleaner way of handling sub/"virtual" components -- logic for paths feels duplicated from Load()
			sub.logicalPath = path.Join(first.logicalPath, sub.Name)
			sub.physicalPath = first.physicalPath
			rest = append(rest, sub)
			echo(3, fmt.Sprintf(`Added component to queue: "%v"`, sub.logicalPath))
		}

		echo(2, "Executing hook: Before-Install")
		if err := first.beforeInstall(); err != nil {
			echo(3, fmt.Errorf(`error running "before-install" hook: %w`, err))
		}

		installer, err := first.ToInstallable()
		if err != nil {
			return visited, fmt.Errorf(`installing component "%v": %w`, first.logicalPath, err)
		}
		if installer != nil {
			echo(2, fmt.Sprintf("Validating installer coordinate: %+v", installer))
			if err := installer.Validate(); err != nil {
				return visited, fmt.Errorf(`validation failed for component coordinate "%+v": %w`, installer, err)
			}

			echo(2, "Computing installation path")
			installPath, err := installer.GetInstallPath()
			if err != nil {
				return visited, err
			}
			echo(3, fmt.Sprintf(`Installation path: "%v"`, installPath))
			first.physicalPath = installPath // overwite to final/"real" location

			echo(2, "Cleaning previous installation")
			if err := installer.Clean(); err != nil {
				return visited, fmt.Errorf(`error cleaning components %s: %w`, first.logicalPath, err)
			}

			echo(2, "Installing")
			if err := installer.Install(); err != nil {
				return visited, fmt.Errorf(`installing component "%v": %w`, first.Name, err)
			}
			echo(3, fmt.Sprintf(`Installed component to: "%v"`, installPath))

			// add remote components to the queue (i.e subcomponents of kind "component")
			if first.Kind == "" || strings.EqualFold(first.Kind, "component") {
				remoteComponentPath := filepath.Join(installPath, first.Path)
				echo(3, fmt.Sprintf(`Adding fetched remote component to queue: "%v"`, remoteComponentPath))
				echo(4, fmt.Sprintf(`Loading component: "%s"`, remoteComponentPath))
				remoteComponent, err := Load(remoteComponentPath, first.logicalPath)
				if err != nil {
					return visited, fmt.Errorf(`loading component from path "%v": %w`, installPath, err)
				}
				echo(5, fmt.Sprintf(`Loaded component with logical path: "%s"`, remoteComponent.logicalPath))
				rest = append(rest, remoteComponent)
				echo(4, fmt.Sprintf(`Added remote component to queue: "%v"`, remoteComponent.logicalPath))
			}
		}

		echo(2, "Executing hook: After-Install")
		if err := first.afterInstall(); err != nil {
			return visited, fmt.Errorf(`error running "after-install" hook: %w`, err)
		}

		visited = append(visited, first)
		echo(2, "Installation complete")

		// reset queue to rest
		queue = rest
	}
}

// type iterator = func(c Component) error

// func Iterate(startPath string, visit iterator) ([]Component, error) {
// 	// Load starting component
// 	component, err := Load(startPath)
// 	if err != nil {
// 		return nil, fmt.Errorf(`failed to load component at "%v": %w`, startPath, err)
// 	}
// 	// Initialize the tree tracking properties of the component
// 	component.PhysicalPath = startPath
// 	component.LogicalPath = fmt.Sprintf("%v", os.PathSeparator)

// 	return iterate(visit, []Component{component}, []Component{})
// }

// func iterate(visit iterator, queue []Component, visited []Component) ([]Component, error) {
// 	//----------------------------------------------------------------------------
// 	// Base case
// 	if len(queue) == 0 {
// 		return visited, nil
// 	}

// 	//----------------------------------------------------------------------------
// 	// Recursive case
// 	first, rest := queue[0], queue[1:]
// 	environments := []string{"common"}
// 	if err := first.LoadConfig(environments); err != nil {
// 		return visited, fmt.Errorf(`error loading configuration "%v" for component %+v: %w`, environments, first, err)
// 	}

// 	// Visit the component
// 	if err := visit(first); err != nil {
// 		return visited, fmt.Errorf(`error visiting component %+v during component iteration: %w`, first, err)
// 	}

// 	// Add all children to the queue
// 	for _, child := range first.Subcomponents {
// 		childType := strings.ToLower(child.ComponentType)
// 		// If of type `component`, load the remote location on disk and enqueue.
// 		// Else, it is inline, so just enqueue the inlined component
// 		if childType == "" || childType == "component" {
// 			installer, err := child.ToInstallable()
// 			if err != nil {
// 				return visited, fmt.Errorf(`error converting subcomponent %+v to installable: %w`, child, err)
// 			}
// 			physicalPath, err := installer.GetInstallPath()
// 			if err != nil {
// 				return visited, fmt.Errorf(`error computing installation path of subcomponent %+v: %w`, child, err)
// 			}
// 			remote, err := Load(physicalPath)
// 			if err != nil {
// 				return visited, fmt.Errorf(`error loading subcomponent in path "%+v": %w`, physicalPath, err)
// 			}
// 			// Remote components have their own component paths
// 			remote.PhysicalPath = physicalPath
// 			remote.LogicalPath = path.Join(first.LogicalPath, child.Name)
// 			rest = append(rest, remote)
// 		} else {
// 			// Inlined components inherit component paths from the parent
// 			child.PhysicalPath = first.PhysicalPath
// 			child.LogicalPath = path.Join(first.LogicalPath, first.Name)
// 			rest = append(rest, child)
// 		}
// 	}

// 	return iterate(visit, rest, visited)
// }
