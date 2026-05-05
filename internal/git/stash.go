package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type StashEntry struct {
	Index     int
	Message   string
	Author    string
	Timestamp string
}

func GetStashList(repoPath string) ([]StashEntry, error) {
	cmd := exec.Command("git", "stash", "list", "--format=%H|%gd|%s|%an|%at")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git stash list in %s: %w", repoPath, err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var stashes []StashEntry
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}
		index := -1
		stashRef := parts[1]
		stashRef = strings.TrimPrefix(stashRef, "stash@{")
		stashRef = strings.TrimSuffix(stashRef, "}")
		if idx, err := strconv.Atoi(stashRef); err == nil {
			index = idx
		}
		stashes = append(stashes, StashEntry{
			Index:     index,
			Message:   parts[2],
			Author:    parts[3],
			Timestamp: parts[4],
		})
	}
	return stashes, nil
}

func ApplyStash(repoPath string, index int) error {
	cmd := exec.Command("git", "stash", "apply", fmt.Sprintf("stash@{%d}", index))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git stash apply in %s: %w", repoPath, err)
	}
	return nil
}

func DropStash(repoPath string, index int) error {
	cmd := exec.Command("git", "stash", "drop", fmt.Sprintf("stash@{%d}", index))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git stash drop in %s: %w", repoPath, err)
	}
	return nil
}