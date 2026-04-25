// Package api provides HTTP handlers for the Git Dashboard.
package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"
	"sync"
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
	Hash      string    `json:"hash"`
	Diff      string    `json:"diff"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

// ReposResponse is the API response for the lightweight repos listing.
type ReposResponse struct {
	Repos []Repo `json:"repos"`
}

// Repo is a lightweight representation of a repository.
type Repo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// GetCommitsFunc is the signature for retrieving commits from a repo.
type GetCommitsFunc func(repoPath string, n int) ([]Commit, error)

// GetDiffFunc is the signature for retrieving a diff from a repo.
type GetDiffFunc func(repoPath, hash string) (string, error)

// GetCommitInfoFunc is the signature for retrieving commit metadata.
type GetCommitInfoFunc func(repoPath, hash string) (message, author string, timestamp time.Time, err error)

// PullFunc is the signature for pulling a repo.
type PullFunc func(repoPath string) error

// SearchFunc is the signature for searching commits.
type SearchFunc func(query, repoPath, author, dateFrom string) []SearchResult

// SearchResult represents a single search result.
type SearchResult struct {
	RepoPath string `json:"repoPath"`
	Hash     string `json:"hash"`
	Message  string `json:"message"`
	Author   string `json:"author"`
	Score    float64 `json:"score"`
}

// PullStatus represents the status of a pull operation for a repository.
type PullStatus struct {
	Repo         string     `json:"repo"`
	LastPullTime *time.Time `json:"lastPullTime,omitempty"`
	LastError    string     `json:"lastError,omitempty"`
	InProgress   bool       `json:"inProgress"`
}

// PullStatusResponse is the API response for pull status.
type PullStatusResponse struct {
	Statuses []PullStatus `json:"statuses"`
}

// HandlerConfig holds the dependencies for the API handlers.
type HandlerConfig struct {
	Repos             []string
	GetCommitsFunc    GetCommitsFunc
	GetDiffFunc       GetDiffFunc
	GetCommitInfoFunc GetCommitInfoFunc
	PullFunc          PullFunc
	SearchFunc        SearchFunc
}

// Handler holds the dependencies for the API handlers.
type Handler struct {
	repos         []string
	getCommits    GetCommitsFunc
	getDiff       GetDiffFunc
	getCommitInfo GetCommitInfoFunc
	pullRepo      PullFunc
	search        SearchFunc

	pullMu       sync.RWMutex
	inProgress   map[string]bool
	lastPullTime map[string]time.Time
	lastPullErr  map[string]string
}

// NewHandler creates a new Handler from a HandlerConfig.
func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		repos:         cfg.Repos,
		getCommits:    cfg.GetCommitsFunc,
		getDiff:       cfg.GetDiffFunc,
		getCommitInfo: cfg.GetCommitInfoFunc,
		pullRepo:      cfg.PullFunc,
		search:        cfg.SearchFunc,
		inProgress:    make(map[string]bool),
		lastPullTime:  make(map[string]time.Time),
		lastPullErr:   make(map[string]string),
	}
}

// RegisterRoutes registers API routes on mux using the provided config.
// It returns the Handler so the caller can set repos via SetRepos.
func RegisterRoutes(mux *http.ServeMux, cfg HandlerConfig) *Handler {
	h := NewHandler(cfg)
	mux.HandleFunc("/api/projects", h.projects)
	mux.HandleFunc("/api/repos", h.listRepos)
	mux.HandleFunc("/api/diff", h.diff)
	mux.HandleFunc("/api/pull", h.pull)
	mux.HandleFunc("/api/pull/status", h.pullStatus)
	mux.HandleFunc("/api/search", h.searchHandler)
	return h
}

// SetRepos updates the list of repos the handler serves.
func (h *Handler) SetRepos(repos []string) {
	h.repos = repos
}

// repos handles GET /api/repos
func (h *Handler) listRepos(w http.ResponseWriter, r *http.Request) {
	var repoList []Repo
	for _, repoPath := range h.repos {
		repoList = append(repoList, Repo{
			Name: filepath.Base(repoPath),
			Path: repoPath,
		})
	}

	if repoList == nil {
		repoList = []Repo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ReposResponse{Repos: repoList})
}

// projects handles GET /api/projects
func (h *Handler) projects(w http.ResponseWriter, r *http.Request) {
	const commitsPerProject = 10
	var projects []Project

	for _, repoPath := range h.repos {
		commits, err := h.getCommits(repoPath, commitsPerProject)
		if err != nil {
			continue
		}
		apiCommits := make([]Commit, len(commits))
		copy(apiCommits, commits)

		var lastCommitAt time.Time
		if len(commits) > 0 {
			lastCommitAt = commits[0].Timestamp
		}

		projects = append(projects, Project{
			Name:         filepath.Base(repoPath),
			Path:         repoPath,
			LastCommitAt: lastCommitAt,
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

	resp := DiffResponse{Hash: hash, Diff: raw}

	if h.getCommitInfo != nil {
		if msg, author, ts, err := h.getCommitInfo(repo, hash); err == nil {
			resp.Message = msg
			resp.Author = author
			resp.Timestamp = ts
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// pull handles POST /api/pull with JSON body {"path": "/repo/path"}
func (h *Handler) pull(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RepoPath string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RepoPath == "" {
		http.Error(w, "missing path parameter", http.StatusBadRequest)
		return
	}

	valid := false
	for _, repo := range h.repos {
		if repo == req.RepoPath {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "unknown repository", http.StatusNotFound)
		return
	}

	h.pullMu.Lock()
	h.inProgress[req.RepoPath] = true
	h.pullMu.Unlock()

	err := h.pullRepo(req.RepoPath)

	h.pullMu.Lock()
	h.inProgress[req.RepoPath] = false
	if err != nil {
		h.lastPullErr[req.RepoPath] = err.Error()
	} else {
		delete(h.lastPullErr, req.RepoPath)
		h.lastPullTime[req.RepoPath] = time.Now()
	}
	h.pullMu.Unlock()

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// pullStatus handles GET /api/pull/status
func (h *Handler) pullStatus(w http.ResponseWriter, r *http.Request) {
	h.pullMu.RLock()
	defer h.pullMu.RUnlock()

	statuses := make([]PullStatus, 0, len(h.repos))
	for _, repo := range h.repos {
		ps := PullStatus{
			Repo:       repo,
			InProgress: h.inProgress[repo],
		}
		if t, ok := h.lastPullTime[repo]; ok {
			ps.LastPullTime = &t
		}
		if err, ok := h.lastPullErr[repo]; ok {
			ps.LastError = err
		}
		statuses = append(statuses, ps)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PullStatusResponse{Statuses: statuses})
}

// search handles GET /api/search?q=<query>&repo=<path>&author=<name>
func (h *Handler) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	repoPath := r.URL.Query().Get("repo")
	author := r.URL.Query().Get("author")
	dateFrom := r.URL.Query().Get("dateFrom")

	if h.search == nil {
		http.Error(w, "search not configured", http.StatusServiceUnavailable)
		return
	}

	results := h.search(query, repoPath, author, dateFrom)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"results": results})
}
