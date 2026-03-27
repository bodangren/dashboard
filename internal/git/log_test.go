package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// initTestRepo creates a temporary git repo with commits for testing.
func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test Author",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test Author",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_DATE=2024-01-15T10:00:00Z",
			"GIT_COMMITTER_DATE=2024-01-15T10:00:00Z",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test Author")

	// First commit with body
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("hello\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit\n\nThis is the commit body.")

	// Second commit, no body
	f2 := filepath.Join(dir, "main.go")
	if err := os.WriteFile(f2, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "add main.go")

	return dir
}

func TestGetCommits_returnsCommits(t *testing.T) {
	dir := initTestRepo(t)

	commits, err := GetCommits(dir, 10)
	if err != nil {
		t.Fatalf("GetCommits returned error: %v", err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}

	// Most recent commit first
	latest := commits[0]
	if latest.Message != "add main.go" {
		t.Errorf("latest message: got %q, want %q", latest.Message, "add main.go")
	}
	if latest.Body != "" {
		t.Errorf("latest body: got %q, want empty", latest.Body)
	}
	if latest.Author != "Test Author" {
		t.Errorf("latest author: got %q, want %q", latest.Author, "Test Author")
	}
	if len(latest.Hash) != 7 {
		t.Errorf("hash length: got %d, want 7", len(latest.Hash))
	}
	if latest.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}

	// Second commit has body
	second := commits[1]
	if second.Message != "initial commit" {
		t.Errorf("second message: got %q, want %q", second.Message, "initial commit")
	}
	if second.Body != "This is the commit body." {
		t.Errorf("second body: got %q, want %q", second.Body, "This is the commit body.")
	}
}

func TestGetCommits_respectsLimit(t *testing.T) {
	dir := initTestRepo(t)

	commits, err := GetCommits(dir, 1)
	if err != nil {
		t.Fatalf("GetCommits returned error: %v", err)
	}
	if len(commits) != 1 {
		t.Errorf("expected 1 commit with limit=1, got %d", len(commits))
	}
}

func TestGetCommits_invalidRepo(t *testing.T) {
	_, err := GetCommits(t.TempDir(), 10)
	if err == nil {
		t.Error("expected error for non-git directory, got nil")
	}
}

func TestGetCommits_timestampParsed(t *testing.T) {
	dir := initTestRepo(t)

	commits, err := GetCommits(dir, 10)
	if err != nil {
		t.Fatalf("GetCommits returned error: %v", err)
	}
	for _, c := range commits {
		if c.Timestamp.Before(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
			t.Errorf("timestamp looks wrong: %v", c.Timestamp)
		}
	}
}
