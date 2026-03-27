package git

import (
	"fmt"
	"log"
	"os/exec"
)

// PullRepo runs `git pull` in the repository at repoPath.
// Errors are returned so the caller (scheduler) can log them without crashing.
func PullRepo(repoPath string) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("git pull in %s: %w\n%s", repoPath, err, string(out))
		log.Printf("pull error: %v", err)
		return err
	}
	return nil
}
