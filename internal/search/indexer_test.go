package search

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIndexer_BuildFromRepos(t *testing.T) {
	idx := NewIndexer()

	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")

	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit")

	if err := idx.BuildFromRepos([]string{dir}); err != nil {
		t.Fatalf("BuildFromRepos failed: %v", err)
	}

	results := idx.Search("initial")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestIndexer_UpdateRepo(t *testing.T) {
	idx := NewIndexer()

	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")

	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit")

	if err := idx.UpdateRepo(dir); err != nil {
		t.Fatalf("UpdateRepo failed: %v", err)
	}

	results := idx.Search("initial")
	if len(results) != 1 {
		t.Fatalf("expected 1 result after initial update, got %d", len(results))
	}

	f2 := filepath.Join(dir, "file2.txt")
	if err := os.WriteFile(f2, []byte("world"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "add second file")

	if err := idx.UpdateRepo(dir); err != nil {
		t.Fatalf("UpdateRepo failed: %v", err)
	}

	results = idx.Search("second")
	if len(results) != 1 {
		t.Fatalf("expected 1 result after second update, got %d", len(results))
	}
}

func TestIndexer_SearchWithFilters(t *testing.T) {
	idx := NewIndexer()

	results := idx.SearchWithFilters("test", "/nonexistent", "", "")
	if len(results) != 0 {
		t.Errorf("expected 0 results for nonexistent path, got %d", len(results))
	}
}