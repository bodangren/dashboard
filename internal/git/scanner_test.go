package git

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// makeGitRepo creates a directory with a .git subdirectory inside t.TempDir().
func makeGitRepo(t *testing.T, root, name string) string {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0755); err != nil {
		t.Fatalf("makeGitRepo: %v", err)
	}
	return dir
}

func TestScanRepos_findsTopLevelRepos(t *testing.T) {
	root := t.TempDir()
	a := makeGitRepo(t, root, "projectA")
	b := makeGitRepo(t, root, "projectB")

	got, err := ScanRepos(root)
	if err != nil {
		t.Fatalf("ScanRepos returned error: %v", err)
	}
	want := []string{a, b}
	sort.Strings(got)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("got %d repos, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("repo[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestScanRepos_findsNestedRepos(t *testing.T) {
	root := t.TempDir()
	nested := makeGitRepo(t, filepath.Join(root, "parent"), "child")

	got, err := ScanRepos(root)
	if err != nil {
		t.Fatalf("ScanRepos returned error: %v", err)
	}
	if len(got) != 1 || got[0] != nested {
		t.Errorf("got %v, want [%s]", got, nested)
	}
}

func TestScanRepos_ignoresNonGitDirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "notARepo"), 0755); err != nil {
		t.Fatal(err)
	}

	got, err := ScanRepos(root)
	if err != nil {
		t.Fatalf("ScanRepos returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no repos, got %v", got)
	}
}

func TestScanRepos_doesNotDescendIntoGitRepo(t *testing.T) {
	// A .git dir inside a repo should not be treated as a separate repo.
	root := t.TempDir()
	repoDir := makeGitRepo(t, root, "myrepo")
	// Create a nested .git dir inside the repo — should not be picked up.
	if err := os.MkdirAll(filepath.Join(repoDir, "subdir", ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	got, err := ScanRepos(root)
	if err != nil {
		t.Fatalf("ScanRepos returned error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 repo, got %v", got)
	}
}
