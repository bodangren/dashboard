package git

import (
	"fmt"
	"os/exec"
)

// GetDiff returns the unified diff for the commit identified by hash
// in the git repository at repoPath.
func GetDiff(repoPath, hash string) (string, error) {
	cmd := exec.Command("git", "show", "--unified=3", hash)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git show %s in %s: %w", hash, repoPath, err)
	}
	return string(out), nil
}
