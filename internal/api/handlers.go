// Package api provides HTTP handlers for the Git Dashboard.
package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"
	"time"
)

// Commit is the API representation of a single git commit.
type Commit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Body      string    `json:"body"`
	Notes     string    `json:"notes"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

// Project is the API representation of a git repository with its commits.
type Project struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	LastCommitAt time.Time `json:"last_commit_at"`
	Commits      []Commit  `json:"commits"`
}

// DiffResponse is the API response for a single commit diff.
type DiffResponse struct {
	Hash string `json:"hash"`
	Diff string `json:"diff"`
}

// GetCommitsFunc is the signature for retrieving commits from a repo.
type GetCommitsFunc func(repoPath string, n int) ([]Commit, error)

// GetDiffFunc is the signature for retrieving a diff from a repo.
type GetDiffFunc func(repoPath, hash string) (string, error)

// Handler holds the dependencies for the API handlers.
type Handler struct {
	repos      []string
	getCommits GetCommitsFunc
	getDiff    GetDiffFunc
}

// RegisterRoutes registers API routes on mux using the real git functions.
// It returns the Handler so the caller can set repos via SetRepos.
func RegisterRoutes(mux *http.ServeMux) *Handler {
	h := &Handler{
		repos:      nil,
		getCommits: defaultGetCommits,
		getDiff:    defaultGetDiff,
	}
	mux.HandleFunc("/api/projects", h.projects)
	mux.HandleFunc("/api/diff", h.diff)
	return h
}

// SetRepos updates the list of repos the handler serves.
func (h *Handler) SetRepos(repos []string) {
	h.repos = repos
}

// projects handles GET /api/projects
func (h *Handler) projects(w http.ResponseWriter, r *http.Request) {
	const commitsPerProject = 10
	var projects []Project

	for _, repoPath := range h.repos {
		commits, err := h.getCommits(repoPath, commitsPerProject)
		if err != nil || len(commits) == 0 {
			continue
		}
		apiCommits := make([]Commit, len(commits))
		copy(apiCommits, commits)

		projects = append(projects, Project{
			Name:         filepath.Base(repoPath),
			Path:         repoPath,
			LastCommitAt: commits[0].Timestamp,
			Commits:      apiCommits,
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastCommitAt.After(projects[j].LastCommitAt)
	})

	if projects == nil {
		projects = []Project{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// diff handles GET /api/diff?repo=<path>&hash=<hash>
func (h *Handler) diff(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	hash := r.URL.Query().Get("hash")
	if repo == "" || hash == "" {
		http.Error(w, "missing repo or hash parameter", http.StatusBadRequest)
		return
	}

	raw, err := h.getDiff(repo, hash)
	if err != nil {
		http.Error(w, "failed to get diff", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DiffResponse{Hash: hash, Diff: raw})
}

// defaultGetCommits and defaultGetDiff bridge to the git package at runtime.
// They are replaced by the import in main.go to avoid circular imports.
var defaultGetCommits GetCommitsFunc = func(repoPath string, n int) ([]Commit, error) {
	return nil, nil
}

var defaultGetDiff GetDiffFunc = func(repoPath, hash string) (string, error) {
	return "", nil
}

// SetGitFuncs injects real git implementations (called from main).
func SetGitFuncs(commits GetCommitsFunc, diff GetDiffFunc) {
	defaultGetCommits = commits
	defaultGetDiff = diff
}
