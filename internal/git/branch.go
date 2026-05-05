package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git branch in %s: %w", repoPath, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var branches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches, nil
}

func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git branch --show-current in %s: %w", repoPath, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func CreateBranch(repoPath, branchName string) error {
	cmd := exec.Command("git", "branch", branchName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git branch %s in %s: %w", branchName, repoPath, err)
	}
	return nil
}

func DeleteBranch(repoPath, branchName string) error {
	cmd := exec.Command("git", "branch", "-d", branchName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git branch -d %s in %s: %w", branchName, repoPath, err)
	}
	return nil
}

func CheckoutBranch(repoPath, branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git checkout %s in %s: %w", branchName, repoPath, err)
	}
	return nil
}