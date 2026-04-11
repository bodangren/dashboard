package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

// CommitInfo holds metadata about a commit.
type CommitInfo struct {
	Message   string
	Author    string
	Timestamp time.Time
}

// GetCommitInfo returns metadata about a commit (message, author, timestamp).
func GetCommitInfo(repoPath, hash string) (CommitInfo, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%s%x1f%an%x1f%ct", hash)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return CommitInfo{}, fmt.Errorf("git log %s in %s: %w", hash, repoPath, err)
	}

	parts := strings.SplitN(strings.TrimSpace(string(out)), "\x1f", 3)
	if len(parts) < 3 {
		return CommitInfo{}, fmt.Errorf("unexpected git log output format")
	}

	ts, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return CommitInfo{}, fmt.Errorf("parse timestamp: %w", err)
	}

	return CommitInfo{
		Message:   parts[0],
		Author:    parts[1],
		Timestamp: time.Unix(ts, 0).UTC(),
	}, nil
}
