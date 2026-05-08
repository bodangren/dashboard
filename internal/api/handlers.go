// Package api provides HTTP handlers for the Git Dashboard.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	LastCommitAt time.Time   `json:"last_commit_at"`
	Commits      []Commit    `json:"commits"`
	Health       *RepoHealth `json:"health,omitempty"`
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

// CommitSearchFunc is the signature for searching commits with full query support.
type CommitSearchFunc func(q *CommitSearchQuery) []SearchResult

// GetHealthFunc is the signature for retrieving repo health.
type GetHealthFunc func(repoPath string) (RepoHealth, error)

// GetBranchesFunc is the signature for retrieving branches.
type GetBranchesFunc func(repoPath string) ([]string, error)

// CreateBranchFunc is the signature for creating a branch.
type CreateBranchFunc func(repoPath, branch string) error

// DeleteBranchFunc is the signature for deleting a branch.
type DeleteBranchFunc func(repoPath, branch string) error

// CheckoutBranchFunc is the signature for checking out a branch.
type CheckoutBranchFunc func(repoPath, branch string) error

// GetStashFunc is the signature for retrieving stash list.
type GetStashFunc func(repoPath string) ([]StashEntry, error)

// ApplyStashFunc is the signature for applying a stash.
type ApplyStashFunc func(repoPath string, index int) error

// DropStashFunc is the signature for dropping a stash.
type DropStashFunc func(repoPath string, index int) error

// StashEntry represents a git stash entry.
type StashEntry struct {
	Index     int    `json:"index"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	RepoPath  string    `json:"repoPath"`
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
}

// DirtyStatus represents the dirty working tree status of a repo.
type DirtyStatus struct {
	Modified  int `json:"modified"`
	Staged    int `json:"staged"`
	Untracked int `json:"untracked"`
	Total     int `json:"total"`
}

// BranchDivergence represents ahead/behind counts relative to upstream.
type BranchDivergence struct {
	Ahead  int `json:"ahead"`
	Behind int `json:"behind"`
}

// StaleBranchInfo represents stale branch information.
type StaleBranchInfo struct {
	Count    int      `json:"count"`
	Branches []string `json:"branches,omitempty"`
}

// RepoHealth represents the overall health status of a repository.
type RepoHealth struct {
	Dirty         DirtyStatus      `json:"dirty"`
	Divergence    BranchDivergence `json:"divergence"`
	StaleBranches StaleBranchInfo  `json:"staleBranches"`
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

// ExportResponse is the API response for data export.
type ExportResponse struct {
	SchemaVersion string    `json:"schemaVersion"`
	ExportedAt    time.Time `json:"exportedAt"`
	Repos         []Repo    `json:"repos"`
}

// ImportRequest is the API request for data import.
type ImportRequest struct {
	SchemaVersion string         `json:"schemaVersion"`
	Repos         []Repo         `json:"repos"`
	Preferences   map[string]interface{} `json:"preferences"`
	Tags          map[string][]string   `json:"tags"`
}

// ImportResponse is the API response for data import.
type ImportResponse struct {
	Status        string `json:"status"`
	ImportedRepos int    `json:"importedRepos"`
}

// ImportErrorResponse is the API error response for data import.
type ImportErrorResponse struct {
	Error string `json:"error"`
}

// HandlerConfig holds the dependencies for the API handlers.
type HandlerConfig struct {
	Repos              []string
	GetCommitsFunc     GetCommitsFunc
	GetDiffFunc        GetDiffFunc
	GetCommitInfoFunc  GetCommitInfoFunc
	PullFunc           PullFunc
	SearchFunc         SearchFunc
	CommitSearchFunc   CommitSearchFunc
	GetHealthFunc      GetHealthFunc
	GetBranchesFunc    GetBranchesFunc
	CreateBranchFunc   CreateBranchFunc
	DeleteBranchFunc   DeleteBranchFunc
	CheckoutBranchFunc CheckoutBranchFunc
	GetStashFunc       GetStashFunc
	ApplyStashFunc     ApplyStashFunc
	DropStashFunc      DropStashFunc
}

// Handler holds the dependencies for the API handlers.
type Handler struct {
	repos          []string
	getCommits     GetCommitsFunc
	getDiff        GetDiffFunc
	getCommitInfo  GetCommitInfoFunc
	pullRepo       PullFunc
	search         SearchFunc
	commitSearch   CommitSearchFunc
	getHealth      GetHealthFunc
	getBranches    GetBranchesFunc
	createBranch   CreateBranchFunc
	deleteBranch   DeleteBranchFunc
	checkoutBranch CheckoutBranchFunc
	getStash       GetStashFunc
	applyStash     ApplyStashFunc
	dropStash      DropStashFunc

	pullMu       sync.RWMutex
	inProgress   map[string]bool
	lastPullTime map[string]time.Time
	lastPullErr  map[string]string

	healthCache   map[string]*healthCacheEntry
	healthCacheMu sync.RWMutex
}

type healthCacheEntry struct {
	health   RepoHealth
	cachedAt time.Time
}

// cacheTTL is the time-to-live for health cache entries.
const healthCacheTTL = 5 * time.Minute

// NewHandler creates a new Handler from a HandlerConfig.
func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		repos:          cfg.Repos,
		getCommits:     cfg.GetCommitsFunc,
		getDiff:        cfg.GetDiffFunc,
		getCommitInfo:  cfg.GetCommitInfoFunc,
		pullRepo:       cfg.PullFunc,
		search:         cfg.SearchFunc,
		commitSearch:   cfg.CommitSearchFunc,
		getHealth:      cfg.GetHealthFunc,
		getBranches:    cfg.GetBranchesFunc,
		createBranch:   cfg.CreateBranchFunc,
		deleteBranch:   cfg.DeleteBranchFunc,
		checkoutBranch: cfg.CheckoutBranchFunc,
		getStash:       cfg.GetStashFunc,
		applyStash:     cfg.ApplyStashFunc,
		dropStash:      cfg.DropStashFunc,
		inProgress:     make(map[string]bool),
		lastPullTime:   make(map[string]time.Time),
		lastPullErr:    make(map[string]string),
		healthCache:    make(map[string]*healthCacheEntry),
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
	mux.HandleFunc("/api/commits/search", h.commitsSearchHandler)
	mux.HandleFunc("/api/health", h.health)
	mux.HandleFunc("/api/branches", h.branches)
	mux.HandleFunc("/api/stash", h.stash)
	mux.HandleFunc("/api/export", h.export)
	mux.HandleFunc("/api/import", h.importData)
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
	var wg sync.WaitGroup

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

	type healthResult struct {
		repoPath string
		health   *RepoHealth
	}
	healthResults := make(chan healthResult, len(projects))

	for _, p := range projects {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()
			hr := healthResult{repoPath: repoPath}
			if h.getHealth != nil {
				if rh, err := h.getHealth(repoPath); err == nil {
					hr.health = &rh
				}
			}
			healthResults <- hr
		}(p.Path)
	}

	wg.Wait()
	close(healthResults)

	healthMap := make(map[string]*RepoHealth, len(projects))
	for hr := range healthResults {
		healthMap[hr.repoPath] = hr.health
	}

	for i := range projects {
		if h, ok := healthMap[projects[i].Path]; ok {
			projects[i].Health = h
		}
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

type CommitSearchResponse struct {
	Results    []CommitSearchResult `json:"results"`
	Total      int                   `json:"total"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
}

type CommitSearchResult struct {
	RepoPath  string    `json:"repoPath"`
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
}

var (
	errEmptyQuery         = errors.New("query text cannot be empty")
	errInvalidDateFormat  = errors.New("invalid date format, use YYYY-MM-DD or relative like 1d, 7d, 30d")
	errInvalidDateRange   = errors.New("since date must be before until date")
	relativeDateRegex    = regexp.MustCompile(`^(\d+)d$`)
)

type CommitSearchQuery struct {
	Q       string
	Author  string
	Repo    string
	Since   *time.Time
	Until   *time.Time
	Limit   int
	Offset  int
}

func parseRelativeDate(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, nil
	}

	if matches := relativeDateRegex.FindStringSubmatch(input); matches != nil {
		days := 0
		for _, c := range matches[1] {
			if c < '0' || c > '9' {
				return time.Time{}, errInvalidDateFormat
			}
			days = days*10 + int(c-'0')
		}
		return time.Now().AddDate(0, 0, -days), nil
	}

	return time.Parse("2006-01-02", input)
}

func (q *CommitSearchQuery) validate() error {
	if q.Q == "" {
		return errEmptyQuery
	}
	if q.Since != nil && q.Until != nil && q.Since.After(*q.Until) {
		return errInvalidDateRange
	}
	return nil
}

func (q *CommitSearchQuery) setDefaults() {
	if q.Limit <= 0 {
		q.Limit = 50
	}
	if q.Limit > 200 {
		q.Limit = 200
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
}

func (h *Handler) commitsSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	commitQuery := &CommitSearchQuery{
		Q:      q,
		Author: r.URL.Query().Get("author"),
		Repo:   r.URL.Query().Get("repo"),
	}

	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		since, err := parseRelativeDate(sinceStr)
		if err != nil {
			http.Error(w, "invalid since parameter: "+err.Error(), http.StatusBadRequest)
			return
		}
		commitQuery.Since = &since
	}

	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		until, err := parseRelativeDate(untilStr)
		if err != nil {
			http.Error(w, "invalid until parameter: "+err.Error(), http.StatusBadRequest)
			return
		}
		commitQuery.Until = &until
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
		commitQuery.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
			return
		}
		commitQuery.Offset = offset
	}

	if h.commitSearch == nil {
		http.Error(w, "commit search not configured", http.StatusServiceUnavailable)
		return
	}

	if err := commitQuery.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	commitQuery.setDefaults()

	results := h.commitSearch(commitQuery)

	apiResults := make([]CommitSearchResult, len(results))
	for i, r := range results {
		apiResults[i] = CommitSearchResult{
			RepoPath:  r.RepoPath,
			Hash:      r.Hash,
			Message:   r.Message,
			Author:    r.Author,
			Timestamp: r.Timestamp,
			Score:     r.Score,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(CommitSearchResponse{
		Results: apiResults,
		Total:   len(apiResults),
		Limit:   commitQuery.Limit,
		Offset:  commitQuery.Offset,
	})
}

// health handles GET /api/health?repo=<path>
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repoPath := r.URL.Query().Get("repo")
	if repoPath == "" {
		http.Error(w, "missing repo parameter", http.StatusBadRequest)
		return
	}

	valid := false
	for _, repo := range h.repos {
		if repo == repoPath {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "unknown repository", http.StatusNotFound)
		return
	}

	if h.getHealth == nil {
		http.Error(w, "health not configured", http.StatusServiceUnavailable)
		return
	}

	health, err := h.getCachedHealth(repoPath)
	if err != nil {
		http.Error(w, "failed to get health: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (h *Handler) getCachedHealth(repoPath string) (RepoHealth, error) {
	h.healthCacheMu.RLock()
	if entry, ok := h.healthCache[repoPath]; ok && time.Since(entry.cachedAt) < healthCacheTTL {
		h.healthCacheMu.RUnlock()
		return entry.health, nil
	}
	h.healthCacheMu.RUnlock()

	health, err := h.getHealth(repoPath)
	if err != nil {
		return RepoHealth{}, err
	}

	h.healthCacheMu.Lock()
	h.healthCache[repoPath] = &healthCacheEntry{
		health:   health,
		cachedAt: time.Now(),
	}
	h.healthCacheMu.Unlock()

	return health, nil
}

func (h *Handler) branches(w http.ResponseWriter, r *http.Request) {
	repoPath := r.URL.Query().Get("repo")
	if repoPath == "" {
		http.Error(w, "missing repo parameter", http.StatusBadRequest)
		return
	}

	valid := false
	for _, repo := range h.repos {
		if repo == repoPath {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "unknown repository", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if h.getBranches == nil {
			http.Error(w, "branches not configured", http.StatusServiceUnavailable)
			return
		}
		branches, err := h.getBranches(repoPath)
		if err != nil {
			http.Error(w, "failed to get branches: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"branches": branches})

	case http.MethodPost:
		var req struct {
			Action string `json:"action"`
			Branch string `json:"branch"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		switch req.Action {
		case "create":
			if h.createBranch == nil {
				http.Error(w, "create branch not configured", http.StatusServiceUnavailable)
				return
			}
			if err := h.createBranch(repoPath, req.Branch); err != nil {
				http.Error(w, "failed to create branch: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "created", "branch": req.Branch})
		case "checkout":
			if h.checkoutBranch == nil {
				http.Error(w, "checkout not configured", http.StatusServiceUnavailable)
				return
			}
			if err := h.checkoutBranch(repoPath, req.Branch); err != nil {
				http.Error(w, "failed to checkout branch: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "checked out", "branch": req.Branch})
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
		}

	case http.MethodDelete:
		var req struct {
			Branch string `json:"branch"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if h.deleteBranch == nil {
			http.Error(w, "delete branch not configured", http.StatusServiceUnavailable)
			return
		}
		if err := h.deleteBranch(repoPath, req.Branch); err != nil {
			http.Error(w, "failed to delete branch: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted", "branch": req.Branch})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) stash(w http.ResponseWriter, r *http.Request) {
	repoPath := r.URL.Query().Get("repo")
	if repoPath == "" {
		http.Error(w, "missing repo parameter", http.StatusBadRequest)
		return
	}

	valid := false
	for _, repo := range h.repos {
		if repo == repoPath {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "unknown repository", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if h.getStash == nil {
			http.Error(w, "stash not configured", http.StatusServiceUnavailable)
			return
		}
		stashes, err := h.getStash(repoPath)
		if err != nil {
			http.Error(w, "failed to get stash: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"stashes": stashes})

	case http.MethodPost:
		var req struct {
			Action string `json:"action"`
			Index  int    `json:"index"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Action == "apply" {
			if h.applyStash == nil {
				http.Error(w, "apply stash not configured", http.StatusServiceUnavailable)
				return
			}
			if err := h.applyStash(repoPath, req.Index); err != nil {
				http.Error(w, "failed to apply stash: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "applied", "index": req.Index})
		} else {
			http.Error(w, "unknown action", http.StatusBadRequest)
		}

	case http.MethodDelete:
		var req struct {
			Index int `json:"index"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if h.dropStash == nil {
			http.Error(w, "drop stash not configured", http.StatusServiceUnavailable)
			return
		}
		if err := h.dropStash(repoPath, req.Index); err != nil {
			http.Error(w, "failed to drop stash: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "dropped", "index": req.Index})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

const currentSchemaVersion = "1.0"

func (h *Handler) export(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	resp := ExportResponse{
		SchemaVersion: currentSchemaVersion,
		ExportedAt:    time.Now(),
		Repos:         repoList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) importData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.SchemaVersion == "" {
		http.Error(w, "schema version required", http.StatusBadRequest)
		return
	}
	if req.SchemaVersion != currentSchemaVersion {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ImportErrorResponse{Error: "unsupported schema version: " + req.SchemaVersion})
		return
	}
	if req.Repos == nil {
		http.Error(w, "repos field required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ImportResponse{
		Status:        "ok",
		ImportedRepos: len(req.Repos),
	})
}
