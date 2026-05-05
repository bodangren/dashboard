package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetBranches(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)
	execBranchGit(t, tmp, "checkout", "-b", "feature-a")
	execBranchGit(t, tmp, "checkout", "-b", "feature-b")
	execBranchGit(t, tmp, "checkout", "main")

	branches, err := GetBranches(tmp)
	if err != nil {
		t.Fatalf("GetBranches: %v", err)
	}
	found := false
	for _, b := range branches {
		if b == "feature-a" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected feature-a in branches, got %v", branches)
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)

	branch, err := GetCurrentBranch(tmp)
	if err != nil {
		t.Fatalf("GetCurrentBranch: %v", err)
	}
	if branch != "main" {
		t.Errorf("expected main, got %s", branch)
	}

	execBranchGit(t, tmp, "checkout", "-b", "feature-x")
	branch, err = GetCurrentBranch(tmp)
	if err != nil {
		t.Fatalf("GetCurrentBranch after checkout: %v", err)
	}
	if branch != "feature-x" {
		t.Errorf("expected feature-x, got %s", branch)
	}
}

func TestCreateBranch(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)

	err := CreateBranch(tmp, "new-branch")
	if err != nil {
		t.Fatalf("CreateBranch: %v", err)
	}

	branches, err := GetBranches(tmp)
	if err != nil {
		t.Fatalf("GetBranches: %v", err)
	}
	found := false
	for _, b := range branches {
		if b == "new-branch" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected new-branch in branches, got %v", branches)
	}
}

func TestDeleteBranch(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)
	execBranchGit(t, tmp, "checkout", "-b", "to-delete")
	execBranchGit(t, tmp, "checkout", "main")

	err := DeleteBranch(tmp, "to-delete")
	if err != nil {
		t.Fatalf("DeleteBranch: %v", err)
	}

	branches, err := GetBranches(tmp)
	if err != nil {
		t.Fatalf("GetBranches: %v", err)
	}
	for _, b := range branches {
		if b == "to-delete" {
			t.Errorf("to-delete should not be in branches")
		}
	}
}

func TestCheckoutBranch(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)

	execBranchGit(t, tmp, "checkout", "-b", "feature-1")
	branchWriteFile(t, tmp, "file.txt", "change on feature")
	execBranchGit(t, tmp, "add", ".")
	execBranchGit(t, tmp, "commit", "-m", "change on feature")

	execBranchGit(t, tmp, "checkout", "main")

	err := CheckoutBranch(tmp, "feature-1")
	if err != nil {
		t.Fatalf("CheckoutBranch: %v", err)
	}

	current, err := GetCurrentBranch(tmp)
	if err != nil {
		t.Fatalf("GetCurrentBranch: %v", err)
	}
	if current != "feature-1" {
		t.Errorf("expected feature-1, got %s", current)
	}
}

func setupTestRepo(t *testing.T, dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	execBranchGit(t, dir, "init")
	execBranchGit(t, dir, "config", "user.email", "test@test.com")
	execBranchGit(t, dir, "config", "user.name", "Test")
	execBranchGit(t, dir, "branch", "-m", "master", "main")
	branchWriteFile(t, dir, "file.txt", "initial content")
	execBranchGit(t, dir, "add", ".")
	execBranchGit(t, dir, "commit", "-m", "initial commit")
}

func execBranchGit(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v in %s: %v", args, dir, err)
	}
}

func branchWriteFile(t *testing.T, dir, name, content string) {
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}