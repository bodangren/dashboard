// Package git provides utilities for interacting with git repositories.
package git

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanRepos walks root up to 2 directory levels deep and returns the absolute
// paths of all directories that contain a .git subdirectory.
// It does not descend into a discovered git repository.
func ScanRepos(root string) ([]string, error) {
	var repos []string

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		if isGitRepo(dir) {
			repos = append(repos, dir)
			continue // do not descend into a git repo
		}
		// One level deeper
		subEntries, err := os.ReadDir(dir)
		if err != nil {
			continue // skip unreadable dirs
		}
		for _, sub := range subEntries {
			if !sub.IsDir() || strings.HasPrefix(sub.Name(), ".") {
				continue
			}
			subDir := filepath.Join(dir, sub.Name())
			if isGitRepo(subDir) {
				repos = append(repos, subDir)
			}
		}
	}

	return repos, nil
}

// isGitRepo reports whether dir contains a .git subdirectory.
func isGitRepo(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && info.IsDir()
}
