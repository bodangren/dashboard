package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dashboard/internal/agents"
)

const testCrontab = `SHELL=/bin/bash
# my agent
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md >> /log/a.log 2>&1
# disabled
# 0 8 * * * cd /home/user/proj2 && gemini -m gemini-2.0 run d.md >> /log/b.log 2>&1
30 2 * * * /usr/bin/cleanup.sh
`

func newAgentHandler(readFn agents.ReadFunc) http.Handler {
	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn)
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)
	return mux
}

func TestAgentsHandler_returnsJSON(t *testing.T) {
	h := newAgentHandler(func() (string, error) { return testCrontab, nil })
	req := httptest.NewRequest("GET", "/api/agents", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}

	var resp AgentsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(resp.Agents))
	}
	if resp.Agents[0].Harness != "opencode" {
		t.Errorf("first agent harness: got %q", resp.Agents[0].Harness)
	}
	if resp.Agents[0].Enabled != true {
		t.Error("first agent should be enabled")
	}
	if resp.Agents[0].SectionHeader != "my agent" {
		t.Errorf("first agent section header: got %q", resp.Agents[0].SectionHeader)
	}
	if resp.Agents[1].Enabled != false {
		t.Error("second agent should be disabled")
	}
}

func TestAgentCreateHandler(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	body, _ := json.Marshal(AgentCreateRequest{
		Schedule:      "0 6 * * *",
		Directory:     "/home/user/new",
		Harness:       "opencode",
		Model:         "zai-coding-plan/glm-5.1",
		Prompt:        "conductor/autonomous_prompt.md",
		LogPath:       "/log/new.log",
		SectionHeader: "New Project",
	})
	req := httptest.NewRequest("POST", "/api/agents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body: %s", rec.Code, rec.Body.String())
	}
	if written == "" {
		t.Fatal("crontab should have been written")
	}
	if !contains(written, "opencode") {
		t.Error("new agent should appear in written crontab")
	}
	if !contains(written, "# /home/user/new") {
		t.Error("section header with directory path should appear in written crontab")
	}
	if !contains(written, "SHELL=/bin/bash") {
		t.Error("existing content should be preserved")
	}
}

func TestAgentDeleteHandler(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	req := httptest.NewRequest("DELETE", "/api/agents/0", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}
	if contains(written, "opencode") {
		t.Error("deleted agent should not appear")
	}
	if !contains(written, "SHELL=/bin/bash") {
		t.Error("non-agent lines should be preserved")
	}
}

func TestAgentToggleHandler(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	req := httptest.NewRequest("PATCH", "/api/agents/0/toggle", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}
	if !contains(written, "# 0 */4") {
		t.Error("toggled agent should now be commented out")
	}
}

func TestToggleAgentHandler_ReturnsUpdatedState(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	req := httptest.NewRequest("PATCH", "/api/agents/0/toggle", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}

	if !contains(written, "# 0 */4") {
		t.Error("toggled agent should now be commented out in crontab")
	}

	var resp AgentJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Enabled {
		t.Error("toggled agent should now be disabled (enabled=false)")
	}
	if resp.Harness != "opencode" {
		t.Errorf("harness should be opencode, got %q", resp.Harness)
	}
	if resp.Schedule != "0 */4 * * *" {
		t.Errorf("schedule should be preserved, got %q", resp.Schedule)
	}
}

func TestAgentLogHandler(t *testing.T) {
	readFn := func() (string, error) { return testCrontab, nil }
	readFileFn := func(path string, n int) (*agents.LogInfo, error) {
		return &agents.LogInfo{Exists: true, Lines: []string{"line1", "line2"}}, nil
	}

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithLogReader(readFileFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	req := httptest.NewRequest("GET", "/api/agents/0/log", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}
	var logInfo agents.LogInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &logInfo); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !logInfo.Exists {
		t.Error("log should exist")
	}
	if len(logInfo.Lines) != 2 {
		t.Errorf("expected 2 log lines, got %d", len(logInfo.Lines))
	}
}

func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}

func TestDiscoverModels_WithExplicitPath(t *testing.T) {
	models := discoverModels("/nonexistent/opencode")
	if models != nil {
		t.Errorf("expected nil for nonexistent binary, got %v", models)
	}
}

func TestDiscoverModels_ParseOutput(t *testing.T) {
	result := parseModelsOutput("model-a\nmodel-b\n\nmodel-c\n")
	if len(result) != 3 {
		t.Fatalf("expected 3 models, got %d", len(result))
	}
	if result[0] != "model-a" || result[1] != "model-b" || result[2] != "model-c" {
		t.Errorf("unexpected models: %v", result)
	}
}

func TestHandleModels_WithBinaryPath(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return "", nil }, WithOpenCodeBinary("/nonexistent/opencode"))
	mux.HandleFunc("/api/models", ah.HandleModels)

	req := httptest.NewRequest("GET", "/api/models", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string][]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["models"] != nil {
		t.Errorf("expected nil models for nonexistent binary, got %v", resp["models"])
	}
}

func TestAgentCreateHandlerUsesRepos(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	ah.SetRepos([]string{"/home/user/proj", "/home/user/new", "/home/user/other"})
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	body, _ := json.Marshal(AgentCreateRequest{
		Schedule:  "0 6 * * *",
		Directory: "/home/user/new",
		Harness:   "opencode",
		Model:     "gpt-4o",
	})
	req := httptest.NewRequest("POST", "/api/agents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body: %s", rec.Code, rec.Body.String())
	}

	// All repos should have section headers
	if !contains(written, "# /home/user/proj") {
		t.Error("repos list should create section headers for all projects")
	}
	if !contains(written, "# /home/user/other") {
		t.Error("empty project should still have section header")
	}
	if !contains(written, "# /home/user/new") {
		t.Error("project with new agent should have section header")
	}
}
