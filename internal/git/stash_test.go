package git

import (
	"testing"
)

func TestGetStashList(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)

	branchWriteFile(t, tmp, "file.txt", "stashed change 1")
	execBranchGit(t, tmp, "add", ".")
	execBranchGit(t, tmp, "stash", "-m", "WIP: first stash")

	branchWriteFile(t, tmp, "file.txt", "stashed change 2")
	execBranchGit(t, tmp, "add", ".")
	execBranchGit(t, tmp, "stash", "-m", "WIP: second stash")

	stashes, err := GetStashList(tmp)
	if err != nil {
		t.Fatalf("GetStashList: %v", err)
	}
	if len(stashes) < 2 {
		t.Errorf("expected at least 2 stashes, got %d", len(stashes))
	}
}

func TestApplyStash(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)

	branchWriteFile(t, tmp, "file.txt", "stashed content")
	execBranchGit(t, tmp, "add", ".")
	execBranchGit(t, tmp, "stash", "-m", "to apply")

	stashes, err := GetStashList(tmp)
	if err != nil {
		t.Fatalf("GetStashList: %v", err)
	}
	if len(stashes) == 0 {
		t.Fatal("expected at least 1 stash")
	}

	err = ApplyStash(tmp, stashes[0].Index)
	if err != nil {
		t.Fatalf("ApplyStash: %v", err)
	}
}

func TestDropStash(t *testing.T) {
	tmp := t.TempDir()
	setupTestRepo(t, tmp)

	branchWriteFile(t, tmp, "file.txt", "to drop")
	execBranchGit(t, tmp, "add", ".")
	execBranchGit(t, tmp, "stash", "-m", "to drop")

	stashes, err := GetStashList(tmp)
	if err != nil {
		t.Fatalf("GetStashList: %v", err)
	}
	if len(stashes) == 0 {
		t.Fatal("expected at least 1 stash")
	}

	err = DropStash(tmp, stashes[0].Index)
	if err != nil {
		t.Fatalf("DropStash: %v", err)
	}

	stashesAfter, err := GetStashList(tmp)
	if err != nil {
		t.Fatalf("GetStashList after drop: %v", err)
	}
	if len(stashesAfter) != 0 {
		t.Errorf("expected 0 stashes after drop, got %d", len(stashesAfter))
	}
}