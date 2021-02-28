package helm

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// BuildInfo encapsulates the parsed output of running `helm template`.
type BuildInfo struct {
	Version      string
	GitCommit    string
	GitTreeState string
	GoVersion    string
}

// Version runs `helm version` and parses the output.
func Version() (v BuildInfo, err error) {
	// Run `helm version` and capture the output
	cmd := exec.Command("helm", "version")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return v, fmt.Errorf(`running %s: %s: %w`, cmd, stderr.String(), err)
	}
	if stderr.String() != "" {
		return v, fmt.Errorf(`running %s: %s`, cmd, stderr.String())
	}

	// capture against stdout
	rgx := regexp.MustCompile(`(?i)Version:"(?P<Version>v\d+\.\d+\.\d+)".*GitCommit:"(?P<GitCommit>[^"]+)".*GitTreeState:"(?P<GitTreeState>[^"]+)".*GoVersion:"(?P<GoVersion>[^"]+)"`)
	matchNames := rgx.SubexpNames()
	for idx, matchValue := range rgx.FindStringSubmatch(stdout.String()) {
		switch matchNames[idx] {
		case "Version":
			v.Version = matchValue
		case "GitCommit":
			v.GitCommit = matchValue
		case "GitTreeState":
			v.GitTreeState = matchValue
		case "GoVersion":
			v.GoVersion = matchValue
		}
	}

	return v, err
}

// IsHelm3 determines if the provided Helm version is major version 3.
func (v BuildInfo) IsHelm3() bool {
	return strings.HasPrefix(strings.ToLower(v.Version), "v3.")
}
