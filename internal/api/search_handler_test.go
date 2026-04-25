package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchHandler(t *testing.T) {
	mockSearch := func(query, repoPath, author, dateFrom string) []SearchResult {
		if query == "test" {
			return []SearchResult{
				{RepoPath: "/tmp/repo1", Hash: "abc1234", Message: "test commit", Author: "Alice", Score: 2.0},
				{RepoPath: "/tmp/repo2", Hash: "def5678", Message: "testing search", Author: "Bob", Score: 1.5},
			}
		}
		return []SearchResult{}
	}

	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos:      []string{"/tmp/repo1", "/tmp/repo2"},
		SearchFunc: mockSearch,
	})
	mux.HandleFunc("/api/search", h.searchHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		Results []SearchResult `json:"results"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(resp.Results))
	}
}

func TestSearchHandler_MissingQuery(t *testing.T) {
	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos:      []string{"/tmp/repo1"},
		SearchFunc: func(q, r, a, d string) []SearchResult { return nil },
	})
	mux.HandleFunc("/api/search", h.searchHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSearchHandler_NoSearchFunc(t *testing.T) {
	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos: []string{"/tmp/repo1"},
	})
	mux.HandleFunc("/api/search", h.searchHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}
}

func TestSearchHandler_WithFilters(t *testing.T) {
	mockSearch := func(query, repoPath, author, dateFrom string) []SearchResult {
		if query == "fix" && repoPath == "/tmp/repo1" && author == "Alice" {
			return []SearchResult{
				{RepoPath: "/tmp/repo1", Hash: "abc1234", Message: "fix auth", Author: "Alice", Score: 2.0},
			}
		}
		return []SearchResult{}
	}

	mux := http.NewServeMux()
	h := NewHandler(HandlerConfig{
		Repos:      []string{"/tmp/repo1", "/tmp/repo2"},
		SearchFunc: mockSearch,
	})
	mux.HandleFunc("/api/search", h.searchHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=fix&repo=/tmp/repo1&author=Alice", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		Results []SearchResult `json:"results"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Author != "Alice" {
		t.Errorf("expected author Alice, got %s", resp.Results[0].Author)
	}
}