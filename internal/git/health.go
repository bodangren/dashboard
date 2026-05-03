package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type DirtyStatus struct {
	Modified   int `json:"modified"`
	Staged     int `json:"staged"`
	Untracked  int `json:"untracked"`
	Total      int `json:"total"`
}

type BranchDivergence struct {
	Ahead  int `json:"ahead"`
	Behind int `json:"behind"`
}

type StaleBranchInfo struct {
	Count    int       `json:"count"`
	Branches []string  `json:"branches,omitempty"`
}

type RepoHealth struct {
	Dirty         DirtyStatus       `json:"dirty"`
	Divergence    BranchDivergence  `json:"divergence"`
	StaleBranches StaleBranchInfo   `json:"staleBranches"`
}

func GetDirtyStatus(repoPath string) (DirtyStatus, error) {
	status := DirtyStatus{}
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return status, fmt.Errorf("git status in %s: %w", repoPath, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if len(line) < 2 {
			continue
		}
		indexStat := line[:2]
		if indexStat[0] != ' ' && indexStat[0] != '?' {
			status.Staged++
		}
		if indexStat[1] == 'M' || indexStat[1] == 'D' {
			status.Modified++
		}
		if line[1] == '?' && line[0] == '?' {
			status.Untracked++
		}
	}
	status.Total = status.Modified + status.Staged + status.Untracked
	return status, nil
}

func GetBranchDivergence(repoPath string) (BranchDivergence, error) {
	div := BranchDivergence{}
	cmd := exec.Command("git", "rev-list", "--count", "--left-right", "@{upstream}...HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return div, nil
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) >= 2 {
		if ahead, err := strconv.Atoi(parts[0]); err == nil {
			div.Ahead = ahead
		}
		if behind, err := strconv.Atoi(parts[1]); err == nil {
			div.Behind = behind
		}
	}
	return div, nil
}

func GetStaleBranches(repoPath string, threshold time.Duration) (StaleBranchInfo, error) {
	info := StaleBranchInfo{}
	cmd := exec.Command("git", "for-each-ref", "--sort=-committerdate", "--format=%(refname:short)|%(committerdate:unix)", "refs/heads/")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return info, fmt.Errorf("git for-each-ref in %s: %w", repoPath, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	cutoff := time.Now().Add(-threshold)
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) < 2 {
			continue
		}
		branch := strings.TrimSpace(parts[0])
		ts, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}
		if time.Unix(ts, 0).Before(cutoff) {
			info.Count++
			if len(info.Branches) < 5 {
				info.Branches = append(info.Branches, branch)
			}
		}
	}
	if len(info.Branches) == 0 {
		info.Branches = nil
	}
	return info, nil
}

func GetRepoHealth(repoPath string) (RepoHealth, error) {
	health := RepoHealth{}
	dirty, err := GetDirtyStatus(repoPath)
	if err == nil {
		health.Dirty = dirty
	}
	div, err := GetBranchDivergence(repoPath)
	if err == nil {
		health.Divergence = div
	}
	stale, err := GetStaleBranches(repoPath, 30*24*time.Hour)
	if err == nil {
		health.StaleBranches = stale
	}
	return health, nil
}