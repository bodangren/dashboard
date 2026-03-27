package git

import (
	"strings"
	"testing"
)

func TestGetDiff_returnsDiff(t *testing.T) {
	dir := initTestRepo(t) // reuses helper from log_test.go

	commits, err := GetCommits(dir, 1)
	if err != nil || len(commits) == 0 {
		t.Fatalf("GetCommits setup failed: %v", err)
	}
	hash := commits[0].Hash

	diff, err := GetDiff(dir, hash)
	if err != nil {
		t.Fatalf("GetDiff returned error: %v", err)
	}
	if diff == "" {
		t.Error("expected non-empty diff")
	}
	if !strings.Contains(diff, "diff --git") {
		t.Errorf("diff missing 'diff --git' header, got: %s", diff[:min(len(diff), 200)])
	}
}

func TestGetDiff_invalidHash(t *testing.T) {
	dir := initTestRepo(t)

	_, err := GetDiff(dir, "0000000")
	if err == nil {
		t.Error("expected error for invalid hash, got nil")
	}
}

func TestGetDiff_invalidRepo(t *testing.T) {
	_, err := GetDiff(t.TempDir(), "abc1234")
	if err == nil {
		t.Error("expected error for non-git dir, got nil")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
