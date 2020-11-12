package url

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// ToPath converts a url to a path like string.
// The path will have OS specific separators.
func ToPath(u string) (string, error) {
	noProtocol, err := removeProtocol(u)
	if err != nil {
		return "", fmt.Errorf(`converting URL "%s" to path: %w`, u, err)
	}

	var pathSegments []string
	for _, v := range strings.Split(noProtocol, "/") {
		if v != "" {
			pathSegments = append(pathSegments, v)
		}
	}

	return filepath.Join(pathSegments...), nil
}

func removeProtocol(repoURL string) (string, error) {
	// Return the original URL if it does not start with a protocol
	if !strings.Contains(repoURL, "://") {
		return repoURL, nil
	}

	// Parse the URL, remove the Scheme and leading "/"
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf(`parsing URL "%s": %w`, repoURL, err)
	}
	u.Scheme = ""

	return strings.TrimLeft(u.String(), "/"), nil
}
