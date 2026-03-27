package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// --- helpers ---

func newTestHandler(repos []string, commitsFn GetCommitsFunc, diffFn GetDiffFunc) http.Handler {
	mux := http.NewServeMux()
	h := &Handler{
		repos:      repos,
		getCommits: commitsFn,
		getDiff:    diffFn,
	}
	mux.HandleFunc("/api/projects", h.projects)
	mux.HandleFunc("/api/diff", h.diff)
	return mux
}

// --- /api/projects tests ---

func TestProjectsHandler_returnsJSON(t *testing.T) {
	t0 := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	t1 := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)

	repos := []string{"/repos/alpha", "/repos/beta"}
	commitsFn := func(repoPath string, n int) ([]Commit, error) {
		if repoPath == "/repos/alpha" {
			return []Commit{{Hash: "aaa1111", Message: "fix bug", Author: "Alice", Timestamp: t0}}, nil
		}
		return []Commit{{Hash: "bbb2222", Message: "add feature", Author: "Bob", Timestamp: t1}}, nil
	}

	h := newTestHandler(repos, commitsFn, nil)
	req := httptest.NewRequest("GET", "/api/projects", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type: got %q, want application/json", ct)
	}

	var projects []Project
	if err := json.Unmarshal(rec.Body.Bytes(), &projects); err != nil {
		t.Fatalf("unmarshal failed: %v\nbody: %s", err, rec.Body.String())
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
}

func TestProjectsHandler_sortedByLatestCommit(t *testing.T) {
	older := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	repos := []string{"/repos/old", "/repos/new"}
	commitsFn := func(repoPath string, n int) ([]Commit, error) {
		if repoPath == "/repos/old" {
			return []Commit{{Hash: "aaa1111", Message: "old commit", Timestamp: older}}, nil
		}
		return []Commit{{Hash: "bbb2222", Message: "new commit", Timestamp: newer}}, nil
	}

	h := newTestHandler(repos, commitsFn, nil)
	req := httptest.NewRequest("GET", "/api/projects", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var projects []Project
	if err := json.Unmarshal(rec.Body.Bytes(), &projects); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if projects[0].Name != "new" {
		t.Errorf("first project should be 'new' (most recent), got %q", projects[0].Name)
	}
}

func TestProjectsHandler_emptyRepos(t *testing.T) {
	h := newTestHandler(nil, func(string, int) ([]Commit, error) { return nil, nil }, nil)
	req := httptest.NewRequest("GET", "/api/projects", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var projects []Project
	if err := json.Unmarshal(rec.Body.Bytes(), &projects); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("expected empty array, got %v", projects)
	}
}

// --- /api/diff tests ---

func TestDiffHandler_returnsDiff(t *testing.T) {
	diffFn := func(repoPath, hash string) (string, error) {
		return "diff --git a/foo.go b/foo.go\n+added line\n", nil
	}

	h := newTestHandler(nil, nil, diffFn)
	req := httptest.NewRequest("GET", "/api/diff?repo=/repos/x&hash=abc1234", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp DiffResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v\nbody: %s", err, rec.Body.String())
	}
	if resp.Hash != "abc1234" {
		t.Errorf("hash: got %q, want abc1234", resp.Hash)
	}
	if resp.Diff == "" {
		t.Error("diff should not be empty")
	}
}

func TestDiffHandler_gitError(t *testing.T) {
	diffFn := func(string, string) (string, error) {
		return "", &testErr{"git show failed"}
	}
	h := newTestHandler(nil, nil, diffFn)
	req := httptest.NewRequest("GET", "/api/diff?repo=/x&hash=abc1234", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

func TestRegisterRoutes_registersEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	for _, path := range []string{"/api/projects", "/api/diff"} {
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Errorf("route %s not registered", path)
		}
	}
}

func TestSetRepos_updatesRepos(t *testing.T) {
	h := &Handler{getCommits: func(string, int) ([]Commit, error) { return nil, nil }}
	h.SetRepos([]string{"/a", "/b"})
	if len(h.repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(h.repos))
	}
}

func TestSetGitFuncs_replacesFuncs(t *testing.T) {
	called := false
	SetGitFuncs(
		func(string, int) ([]Commit, error) { called = true; return nil, nil },
		func(string, string) (string, error) { return "", nil },
	)
	defaultGetCommits("", 0)
	if !called {
		t.Error("SetGitFuncs did not replace defaultGetCommits")
	}
	// Restore no-op defaults
	SetGitFuncs(
		func(string, int) ([]Commit, error) { return nil, nil },
		func(string, string) (string, error) { return "", nil },
	)
}

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }

func TestDiffHandler_missingParams(t *testing.T) {
	h := newTestHandler(nil, nil, func(string, string) (string, error) { return "", nil })

	tests := []struct {
		url  string
		name string
	}{
		{"/api/diff", "no params"},
		{"/api/diff?repo=/x", "missing hash"},
		{"/api/diff?hash=abc1234", "missing repo"},
	}
	for _, tc := range tests {
		req := httptest.NewRequest("GET", tc.url, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: expected 400, got %d", tc.name, rec.Code)
		}
	}
}
