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
0 */4 * * * cd /home/user/proj && opencode --model gpt-4o --prompt t.md >> /log/a.log 2>&1
# disabled
# 0 8 * * * cd /home/user/proj2 && gemini --model gemini-2.0 --prompt d.md >> /log/b.log 2>&1
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
		Schedule:  "0 6 * * *",
		Directory: "/home/user/new",
		Harness:   "codex",
		Model:     "codex-1",
		Prompt:    "fix.md",
		LogPath:   "/log/new.log",
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
	if !contains(written, "codex") {
		t.Error("new agent should appear in written crontab")
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
