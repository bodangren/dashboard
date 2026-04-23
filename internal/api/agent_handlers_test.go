package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("DELETE", "/api/agents/"+agentID, nil)
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

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("PATCH", "/api/agents/"+agentID+"/toggle", nil)
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

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("PATCH", "/api/agents/"+agentID+"/toggle", nil)
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

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("GET", "/api/agents/"+agentID+"/log", nil)
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

func TestAgentDeleteByID(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("DELETE", "/api/agents/"+agentID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]bool
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !resp["ok"] {
		t.Error("expected ok=true in response")
	}

	if contains(written, "gpt-4o") {
		t.Error("deleted agent should not appear in written crontab")
	}

	if !contains(written, "gemini-2.0") {
		t.Error("second agent should still be present after deleting first agent")
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

func TestAgentUpdateHandler(t *testing.T) {
	var written string
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { written = content; return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	body, _ := json.Marshal(AgentCreateRequest{
		Schedule:  "0 6 * * *",
		Directory: "/home/user/proj",
		Harness:   "opencode",
		Model:     "gpt-5",
		Prompt:    "updated.md",
		LogPath:   "/log/updated.log",
	})
	req := httptest.NewRequest("PUT", "/api/agents/"+agentID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}

	if !contains(written, "gpt-5") {
		t.Error("updated agent should have new model")
	}
	if contains(written, "gpt-4o") {
		t.Error("old model should not appear after update")
	}
	if !contains(written, "updated.md") {
		t.Error("updated agent should have new prompt")
	}
	if !contains(written, "/log/updated.log") {
		t.Error("updated agent should have new log path")
	}
}

func TestAgentUpdateHandler_NotFound(t *testing.T) {
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("nonexistent:schedule:model")
	body, _ := json.Marshal(AgentCreateRequest{
		Schedule:  "0 6 * * *",
		Directory: "/home/user/proj",
		Harness:   "opencode",
		Model:     "gpt-5",
	})
	req := httptest.NewRequest("PUT", "/api/agents/"+agentID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for non-existent agent, got %d", rec.Code)
	}
}

func TestAgentUpdateHandler_InvalidJSON(t *testing.T) {
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("PUT", "/api/agents/"+agentID, bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestAgentLogHandler_NoLogPath(t *testing.T) {
	input := `SHELL=/bin/bash
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md
`
	readFn := func() (string, error) { return input, nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn)
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("GET", "/api/agents/"+agentID+"/log", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for agent without log path, got %d", rec.Code)
	}

	var logInfo agents.LogInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &logInfo); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if logInfo.Exists {
		t.Error("log should not exist for agent without log path")
	}
}

func TestAgentsHandler_ListEmpty(t *testing.T) {
	readFn := func() (string, error) { return "SHELL=/bin/bash\n", nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn)
	mux.HandleFunc("/api/agents", ah.HandleAgents)

	req := httptest.NewRequest("GET", "/api/agents", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp AgentsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Agents) != 0 {
		t.Errorf("expected 0 agents for crontab without agents, got %d", len(resp.Agents))
	}
}

func TestAgentCreateHandler_MissingRequiredFields(t *testing.T) {
	readFn := func() (string, error) { return testCrontab, nil }
	writeFn := func(content string) error { return nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithWriteFunc(writeFn))
	mux.HandleFunc("/api/agents", ah.HandleAgents)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	body, _ := json.Marshal(AgentCreateRequest{
		Schedule:  "", // missing
		Directory: "/home/user/new",
		Harness:   "opencode",
	})
	req := httptest.NewRequest("POST", "/api/agents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing schedule, got %d", rec.Code)
	}
}

func TestHandleAgents_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents", ah.HandleAgents)

	req := httptest.NewRequest("DELETE", "/api/agents", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandleAgentAction_MissingAgentID(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	req := httptest.NewRequest("DELETE", "/api/agents/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing agent ID, got %d", rec.Code)
	}
}

func TestHandleAgentAction_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("POST", "/api/agents/"+agentID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestToggleAgent_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("POST", "/api/agents/"+agentID+"/toggle", nil) // POST instead of PATCH
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for wrong method on toggle, got %d", rec.Code)
	}
}

func TestGetLog_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("POST", "/api/agents/"+agentID+"/log", nil) // POST instead of GET
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for wrong method on log, got %d", rec.Code)
	}
}

func TestGetLog_ReadLogError(t *testing.T) {
	readFn := func() (string, error) { return testCrontab, nil }
	readFileFn := func(path string, n int) (*agents.LogInfo, error) {
		return nil, fmt.Errorf("simulated read error")
	}

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn, WithLogReader(readFileFn))
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("GET", "/api/agents/"+agentID+"/log", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 even on read error, got %d", rec.Code)
	}

	var logInfo agents.LogInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &logInfo); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if logInfo.Exists {
		t.Error("log should not exist when read fails")
	}
}

func TestGetLog_AgentNotFound(t *testing.T) {
	readFn := func() (string, error) { return testCrontab, nil }

	mux := http.NewServeMux()
	ah := NewAgentHandler(readFn)
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("nonexistent:schedule:model")
	req := httptest.NewRequest("GET", "/api/agents/"+agentID+"/log", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for non-existent agent, got %d", rec.Code)
	}
}

func TestTriggerAgent_Success(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("POST", "/api/agents/"+agentID+"/trigger", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["status"] != "triggered" {
		t.Errorf("expected status 'triggered', got %q", resp["status"])
	}
	if resp["agent_id"] != "0 */4 * * *:/home/user/proj:gpt-4o" {
		t.Errorf("expected agent_id, got %q", resp["agent_id"])
	}
}

func TestTriggerAgent_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("0 */4 * * *:/home/user/proj:gpt-4o")
	req := httptest.NewRequest("GET", "/api/agents/"+agentID+"/trigger", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestTriggerAgent_AgentNotFound(t *testing.T) {
	mux := http.NewServeMux()
	ah := NewAgentHandler(func() (string, error) { return testCrontab, nil })
	mux.HandleFunc("/api/agents/", ah.HandleAgentAction)

	agentID := url.PathEscape("nonexistent:schedule:model")
	req := httptest.NewRequest("POST", "/api/agents/"+agentID+"/trigger", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body: %s", rec.Code, rec.Body.String())
	}
}
