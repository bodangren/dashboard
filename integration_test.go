package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"dashboard/internal/api"
	gitpkg "dashboard/internal/git"
)

// makeTestRepo creates a temp git repo with one commit and returns its path.
func makeTestRepo(t *testing.T, name string) string {
	t.Helper()
	base := t.TempDir()
	dir := filepath.Join(base, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=t@t.com",
			"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=t@t.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init")
	run("config", "user.email", "t@t.com")
	run("config", "user.name", "Test")
	f := filepath.Join(dir, "readme.txt")
	if err := os.WriteFile(f, []byte("hello\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit")
	return dir
}

func buildTestMux(repos []string) http.Handler {
	// SetGitFuncs must be called before RegisterRoutes so the handler captures
	// the real implementations.
	api.SetGitFuncs(
		func(repoPath string, n int) ([]api.Commit, error) {
			gitCommits, err := gitpkg.GetCommits(repoPath, n)
			if err != nil {
				return nil, err
			}
			out := make([]api.Commit, len(gitCommits))
			for i, c := range gitCommits {
				out[i] = api.Commit{
					Hash: c.Hash, Message: c.Message, Body: c.Body,
					Notes: c.Notes, Author: c.Author, Timestamp: c.Timestamp,
				}
			}
			return out, nil
		},
		func(repoPath, hash string) (string, error) {
			return gitpkg.GetDiff(repoPath, hash)
		},
	)
	mux := http.NewServeMux()
	h := api.RegisterRoutes(mux)
	h.SetRepos(repos)
	return mux
}

func TestServerEndToEnd_projectsEndpoint(t *testing.T) {
	repo := makeTestRepo(t, "myproject")
	mux := buildTestMux([]string{repo})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/api/projects")
	if err != nil {
		t.Fatalf("GET /api/projects: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var projects []struct {
		Name         string    `json:"name"`
		LastCommitAt time.Time `json:"last_commit_at"`
		Commits      []struct {
			Hash    string `json:"hash"`
			Message string `json:"message"`
		} `json:"commits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&projects); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "myproject" {
		t.Errorf("name: got %q, want myproject", projects[0].Name)
	}
	if len(projects[0].Commits) == 0 {
		t.Error("expected at least one commit")
	}
	if projects[0].Commits[0].Message != "initial commit" {
		t.Errorf("commit message: got %q", projects[0].Commits[0].Message)
	}
}

func TestServerEndToEnd_diffEndpoint(t *testing.T) {
	repo := makeTestRepo(t, "diffproject")
	mux := buildTestMux([]string{repo})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// Get the commit hash from /api/projects first
	res, _ := http.Get(srv.URL + "/api/projects")
	var projects []struct {
		Path    string `json:"path"`
		Commits []struct {
			Hash string `json:"hash"`
		} `json:"commits"`
	}
	json.NewDecoder(res.Body).Decode(&projects)
	res.Body.Close()

	hash := projects[0].Commits[0].Hash
	path := projects[0].Path

	diffURL := srv.URL + "/api/diff?repo=" + path + "&hash=" + hash
	res2, err := http.Get(diffURL)
	if err != nil {
		t.Fatalf("GET /api/diff: %v", err)
	}
	defer res2.Body.Close()

	if res2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res2.StatusCode)
	}
	var diffResp struct {
		Hash string `json:"hash"`
		Diff string `json:"diff"`
	}
	if err := json.NewDecoder(res2.Body).Decode(&diffResp); err != nil {
		t.Fatalf("decode diff: %v", err)
	}
	if diffResp.Diff == "" {
		t.Error("expected non-empty diff")
	}
}
