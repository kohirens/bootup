package main

import (
	"path/filepath"
	"strings"
)

// getRepoDir Extract a local dirname from a Git URL.
func getRepoDir(repoLocation, refName string) string {
	if len(repoLocation) < 1 {
		return repoLocation
	}

	baseName := filepath.Base(repoLocation)

	// trim .git from the end
	baseName = strings.TrimRight(baseName, ".git")

	// append ref, branch, or tag
	if refName != "" {
		baseName = baseName + "-" + strings.ReplaceAll(refName, "/", "-")
	}

	return baseName
}
