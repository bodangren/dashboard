package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- helpers ---

func newTestHandler(repos []string, commitsFn GetCommitsFunc, diffFn GetDiffFunc) http.Handler {
	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos:          repos,
		GetCommitsFunc: commitsFn,
		GetDiffFunc:    diffFn,
		PullFunc:       func(string) error { return nil },
	})
	mux.HandleFunc("/api/projects", h.projects)
	mux.HandleFunc("/api/repos", h.listRepos)
	mux.HandleFunc("/api/diff", h.diff)
	mux.HandleFunc("/api/pull", h.pull)
	return mux
}

func newTestHandlerWithPull(repos []string, commitsFn GetCommitsFunc, diffFn GetDiffFunc, pullFn PullFunc) http.Handler {
	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos:          repos,
		GetCommitsFunc: commitsFn,
		GetDiffFunc:    diffFn,
		PullFunc:       pullFn,
	})
	mux.HandleFunc("/api/projects", h.projects)
	mux.HandleFunc("/api/repos", h.listRepos)
	mux.HandleFunc("/api/diff", h.diff)
	mux.HandleFunc("/api/pull", h.pull)
	return mux
}

func newTestHandlerWithCommitInfo(repos []string, diffFn GetDiffFunc, commitInfoFn GetCommitInfoFunc) http.Handler {
	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos:             repos,
		GetCommitInfoFunc: commitInfoFn,
		GetDiffFunc:       diffFn,
		PullFunc:          func(string) error { return nil },
	})
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

func TestProjectsHandler_includesEmptyRepos(t *testing.T) {
	repos := []string{"/repos/empty", "/repos/has-commits"}
	commitsFn := func(repoPath string, n int) ([]Commit, error) {
		if repoPath == "/repos/has-commits" {
			return []Commit{{Hash: "abc1234", Message: "test", Timestamp: time.Now()}}, nil
		}
		return []Commit{}, nil
	}

	h := newTestHandler(repos, commitsFn, nil)
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
	if len(projects) != 2 {
		t.Errorf("expected 2 projects (including empty), got %d", len(projects))
	}
	for _, p := range projects {
		if p.Name == "empty" && len(p.Commits) != 0 {
			t.Errorf("empty repo should have 0 commits, got %d", len(p.Commits))
		}
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

func TestDiffHandler_returnsMetadata(t *testing.T) {
	diffFn := func(repoPath, hash string) (string, error) {
		return "diff --git a/foo.go b/foo.go\n+added line\n", nil
	}
	commitInfoFn := func(repoPath, hash string) (string, string, time.Time, error) {
		return "fix bug", "Alice", time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC), nil
	}

	h := newTestHandlerWithCommitInfo(nil, diffFn, commitInfoFn)
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
	if resp.Message != "fix bug" {
		t.Errorf("message: got %q, want 'fix bug'", resp.Message)
	}
	if resp.Author != "Alice" {
		t.Errorf("author: got %q, want 'Alice'", resp.Author)
	}
	if resp.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestRegisterRoutes_registersEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, HandlerConfig{
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
		GetDiffFunc:    func(string, string) (string, error) { return "", nil },
		PullFunc:       func(string) error { return nil },
	})

	for _, path := range []string{"/api/projects", "/api/diff", "/api/pull"} {
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Errorf("route %s not registered", path)
		}
	}
}

func TestSetRepos_updatesRepos(t *testing.T) {
	h := NewHandler(HandlerConfig{
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	h.SetRepos([]string{"/a", "/b"})
	if len(h.repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(h.repos))
	}
}

func TestNewHandlerConfig_NoGlobals(t *testing.T) {
	repos := []string{"/repos/test"}
	calledCommits := false
	calledDiff := false
	calledPull := false

	cfg := HandlerConfig{
		Repos: repos,
		GetCommitsFunc: func(repoPath string, n int) ([]Commit, error) {
			calledCommits = true
			return []Commit{{Hash: "abc", Message: "test", Timestamp: time.Now()}}, nil
		},
		GetDiffFunc: func(repoPath, hash string) (string, error) {
			calledDiff = true
			return "diff content", nil
		},
		PullFunc: func(repoPath string) error {
			calledPull = true
			return nil
		},
	}

	mux := http.NewServeMux()
	h := NewHandler(cfg)
	mux.HandleFunc("/api/projects", h.projects)
	mux.HandleFunc("/api/diff", h.diff)
	mux.HandleFunc("/api/pull", h.pull)

	req := httptest.NewRequest("GET", "/api/projects", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !calledCommits {
		t.Error("GetCommitsFunc was not called")
	}

	req = httptest.NewRequest("GET", "/api/diff?repo=/repos/test&hash=abc", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if !calledDiff {
		t.Error("GetDiffFunc was not called")
	}

	body := `{"path":"/repos/test"}`
	req = httptest.NewRequest("POST", "/api/pull", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if !calledPull {
		t.Error("PullFunc was not called")
	}
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

// --- /api/pull tests ---

func TestPullHandler_success(t *testing.T) {
	pulled := ""
	pullFn := func(repoPath string) error {
		pulled = repoPath
		return nil
	}

	h := newTestHandlerWithPull([]string{"/repos/alpha"}, nil, nil, pullFn)
	req := httptest.NewRequest("POST", "/api/pull", strings.NewReader(`{"path":"/repos/alpha"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if pulled != "/repos/alpha" {
		t.Errorf("expected pull for /repos/alpha, got %q", pulled)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %q", resp["status"])
	}
}

func TestPullHandler_unknownRepo(t *testing.T) {
	h := newTestHandlerWithPull([]string{"/repos/alpha"}, nil, nil, func(string) error { return nil })
	req := httptest.NewRequest("POST", "/api/pull", strings.NewReader(`{"path":"/repos/unknown"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestPullHandler_missingBody(t *testing.T) {
	h := newTestHandlerWithPull([]string{"/repos/alpha"}, nil, nil, func(string) error { return nil })
	req := httptest.NewRequest("POST", "/api/pull", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestPullHandler_methodNotAllowed(t *testing.T) {
	h := newTestHandlerWithPull([]string{"/repos/alpha"}, nil, nil, func(string) error { return nil })
	req := httptest.NewRequest("GET", "/api/pull", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestPullHandler_pullError(t *testing.T) {
	pullFn := func(string) error {
		return &testErr{"git pull failed: no tracking information"}
	}

	h := newTestHandlerWithPull([]string{"/repos/alpha"}, nil, nil, pullFn)
	req := httptest.NewRequest("POST", "/api/pull", strings.NewReader(`{"path":"/repos/alpha"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["status"] != "error" {
		t.Errorf("expected status 'error', got %q", resp["status"])
	}
	if !strings.Contains(resp["error"], "no tracking information") {
		t.Errorf("error should contain specific message, got %q", resp["error"])
	}
}

func TestReposHandler_returnsOnlyNamesAndPaths(t *testing.T) {
	h := newTestHandler([]string{"/repos/alpha", "/repos/beta"}, nil, nil)
	req := httptest.NewRequest("GET", "/api/repos", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp ReposResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(resp.Repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(resp.Repos))
	}

	if resp.Repos[0].Name != "alpha" || resp.Repos[0].Path != "/repos/alpha" {
		t.Errorf("first repo: got name=%q path=%q, want name=alpha path=/repos/alpha", resp.Repos[0].Name, resp.Repos[0].Path)
	}
	if resp.Repos[1].Name != "beta" || resp.Repos[1].Path != "/repos/beta" {
		t.Errorf("second repo: got name=%q path=%q, want name=beta path=/repos/beta", resp.Repos[1].Name, resp.Repos[1].Path)
	}
}

func TestReposHandler_emptyRepos(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	req := httptest.NewRequest("GET", "/api/repos", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp ReposResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(resp.Repos) != 0 {
		t.Errorf("expected 0 repos, got %d", len(resp.Repos))
	}
}
