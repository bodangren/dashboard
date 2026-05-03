package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestGetDirtyStatus(t *testing.T) {
	tmp := t.TempDir()
	initGitRepo(t, tmp)

	dirty, err := GetDirtyStatus(tmp)
	if err != nil {
		t.Fatalf("GetDirtyStatus on clean repo: %v", err)
	}
	if dirty.Total != 0 {
		t.Errorf("expected 0 total changes, got %d", dirty.Total)
	}

	writeFile(t, tmp, "untracked.txt", "new")

	dirty, err = GetDirtyStatus(tmp)
	if err != nil {
		t.Fatalf("GetDirtyStatus after untracked file: %v", err)
	}
	if dirty.Untracked != 1 {
		t.Errorf("expected 1 untracked, got %d", dirty.Untracked)
	}

	writeFile(t, tmp, "tracked.txt", "content")
	runGit(t, tmp, "add", "tracked.txt")

	dirty, err = GetDirtyStatus(tmp)
	if err != nil {
		t.Fatalf("GetDirtyStatus after staging: %v", err)
	}
	if dirty.Staged != 1 {
		t.Errorf("expected 1 staged, got %d", dirty.Staged)
	}

	runGit(t, tmp, "commit", "-m", "init")

	os.Remove(filepath.Join(tmp, "untracked.txt"))

	dirty, err = GetDirtyStatus(tmp)
	if err != nil {
		t.Fatalf("GetDirtyStatus after commit: %v", err)
	}
	if dirty.Total != 0 {
		t.Errorf("expected 0 total after commit, got %d", dirty.Total)
	}
}

func TestGetBranchDivergence(t *testing.T) {
	tmp := t.TempDir()
	initGitRepo(t, tmp)

	div, err := GetBranchDivergence(tmp)
	if err != nil {
		t.Fatalf("GetBranchDivergence: %v", err)
	}
	if div.Ahead != 0 || div.Behind != 0 {
		t.Errorf("expected 0/0 divergence, got %d/%d", div.Ahead, div.Behind)
	}
}

func TestGetStaleBranches(t *testing.T) {
	tmp := t.TempDir()
	initGitRepo(t, tmp)

	info, err := GetStaleBranches(tmp, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("GetStaleBranches: %v", err)
	}
	if info.Count != 0 {
		t.Errorf("expected 0 stale branches, got %d", info.Count)
	}
}

func TestGetRepoHealth(t *testing.T) {
	tmp := t.TempDir()
	initGitRepo(t, tmp)

	health, err := GetRepoHealth(tmp)
	if err != nil {
		t.Fatalf("GetRepoHealth: %v", err)
	}
	if health.Dirty.Total != 0 {
		t.Errorf("expected clean repo, got %d dirty", health.Dirty.Total)
	}
}

func TestGetDirtyStatusInvalidPath(t *testing.T) {
	_, err := GetDirtyStatus("/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent path")
	}
}

func TestGetStaleBranchesEmpty(t *testing.T) {
	tmp := t.TempDir()
	initGitRepo(t, tmp)

	info, err := GetStaleBranches(tmp, 1*time.Hour)
	if err != nil {
		t.Fatalf("GetStaleBranches: %v", err)
	}
	if info.Count != 0 {
		t.Errorf("expected 0 stale branches in fresh repo, got %d", info.Count)
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
}