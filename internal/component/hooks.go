package component

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ExecuteHook executes the passed hook
func (c *Component) ExecuteHook(hook string, commands []string) (err error) {
	for _, command := range commands {
		if len(command) != 0 {
			cmd := exec.Command("sh", "-c", command)
			cmd.Dir = c.physicalPath
			cmd.Stdout = os.Stdout
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf(`error executing component hook "%v": %w`, hook, err)
			}
		}
	}

	return nil
}

// beforeGenerate executes the 'before-generate' hook (if any) of the component.
func (c *Component) beforeGenerate() (err error) {
	return c.ExecuteHook("before-generate", c.Hooks.BeforeGenerate)
}

// afterGenerate executes the 'after-generate' hook (if any) of the component.
func (c *Component) afterGenerate() (err error) {
	return c.ExecuteHook("after-generate", c.Hooks.AfterGenerate)
}

// beforeInstall executes the 'before-install' hook (if any) of the component.
func (c *Component) beforeInstall() (err error) {
	return c.ExecuteHook("before-install", c.Hooks.BeforeInstall)
}

// afterInstall executes the 'after-install' hook (if any) of the component.
func (c *Component) afterInstall() (err error) {
	return c.ExecuteHook("after-install", c.Hooks.AfterInstall)
}
