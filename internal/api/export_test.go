package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExportHandler_returnsExportData(t *testing.T) {
	repos := []string{"/repos/alpha", "/repos/beta"}
	h := NewHandler(HandlerConfig{
		Repos:          repos,
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/export", h.export)

	req := httptest.NewRequest("GET", "/api/export", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type: got %q, want application/json", ct)
	}

	var resp ExportResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v\nbody: %s", err, rec.Body.String())
	}

	if resp.SchemaVersion != "1.0" {
		t.Errorf("schemaVersion: got %q, want 1.0", resp.SchemaVersion)
	}
	if resp.ExportedAt.IsZero() {
		t.Error("exportedAt should not be zero")
	}
	if len(resp.Repos) != 2 {
		t.Errorf("repos count: got %d, want 2", len(resp.Repos))
	}
}

func TestExportHandler_includesRepoNamesAndPaths(t *testing.T) {
	repos := []string{"/repos/project-a", "/repos/project-b"}
	h := NewHandler(HandlerConfig{
		Repos:          repos,
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/export", h.export)

	req := httptest.NewRequest("GET", "/api/export", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp ExportResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.Repos[0].Name != "project-a" || resp.Repos[0].Path != "/repos/project-a" {
		t.Errorf("first repo: got name=%q path=%q, want name=project-a path=/repos/project-a", resp.Repos[0].Name, resp.Repos[0].Path)
	}
	if resp.Repos[1].Name != "project-b" || resp.Repos[1].Path != "/repos/project-b" {
		t.Errorf("second repo: got name=%q path=%q, want name=project-b path=/repos/project-b", resp.Repos[1].Name, resp.Repos[1].Path)
	}
}

func TestExportHandler_includesTimestamp(t *testing.T) {
	before := time.Now().Add(-time.Second)
	h := NewHandler(HandlerConfig{
		Repos:          []string{},
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/export", h.export)

	req := httptest.NewRequest("GET", "/api/export", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	after := time.Now().Add(time.Second)
	var resp ExportResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.ExportedAt.Before(before) || resp.ExportedAt.After(after) {
		t.Errorf("exportedAt outside expected range: got %v, want between %v and %v", resp.ExportedAt, before, after)
	}
}

func TestExportHandler_methodNotAllowed(t *testing.T) {
	h := NewHandler(HandlerConfig{
		Repos:          []string{},
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/export", h.export)

	req := httptest.NewRequest("POST", "/api/export", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestImportHandler_validatesSchemaVersion(t *testing.T) {
	h := NewHandler(HandlerConfig{
		Repos:          []string{"/repos/test"},
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/import", h.importData)

	body := `{"schemaVersion":"99.0","repos":[],"preferences":{},"tags":{}}`
	req := httptest.NewRequest("POST", "/api/import", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown schema version, got %d", rec.Code)
	}

	var resp ImportErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !strings.Contains(resp.Error, "schema version") {
		t.Errorf("error should mention schema version, got %q", resp.Error)
	}
}

func TestImportHandler_missingRequiredFields(t *testing.T) {
	h := NewHandler(HandlerConfig{
		Repos:          []string{"/repos/test"},
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/import", h.importData)

	tests := []struct {
		name string
		body string
	}{
		{"empty body", `{}`},
		{"missing schemaVersion", `{"repos":[]}`},
		{"missing repos", `{"schemaVersion":"1.0"}`},
	}

	for _, tc := range tests {
		req := httptest.NewRequest("POST", "/api/import", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: expected 400, got %d", tc.name, rec.Code)
		}
	}
}

func TestImportHandler_methodNotAllowed(t *testing.T) {
	h := NewHandler(HandlerConfig{
		Repos:          []string{},
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/import", h.importData)

	req := httptest.NewRequest("GET", "/api/import", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestImportResponse_structure(t *testing.T) {
	h := NewHandler(HandlerConfig{
		Repos:          []string{"/repos/test"},
		GetCommitsFunc: func(string, int) ([]Commit, error) { return nil, nil },
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/import", h.importData)

	body := `{"schemaVersion":"1.0","repos":[{"name":"test","path":"/repos/test"}],"preferences":{},"tags":{}}`
	req := httptest.NewRequest("POST", "/api/import", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp ImportResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("status: got %q, want ok", resp.Status)
	}
	if resp.ImportedRepos != 1 {
		t.Errorf("importedRepos: got %d, want 1", resp.ImportedRepos)
	}
}