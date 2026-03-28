package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"dashboard/internal/agents"
)

type AgentJSON struct {
	LineIndex int    `json:"line_index"`
	Schedule  string `json:"schedule"`
	Directory string `json:"directory"`
	Harness   string `json:"harness"`
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	LogPath   string `json:"log_path"`
	Enabled   bool   `json:"enabled"`
}

type AgentsResponse struct {
	Agents []AgentJSON `json:"agents"`
}

type AgentCreateRequest struct {
	Schedule  string `json:"schedule"`
	Directory string `json:"directory"`
	Harness   string `json:"harness"`
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	LogPath   string `json:"log_path"`
}

type AgentHandler struct {
	readCrontab  agents.ReadFunc
	writeCrontab agents.WriteFunc
	readLog      agents.LogReadFunc
}

type AgentHandlerOption func(*AgentHandler)

func WithWriteFunc(fn agents.WriteFunc) AgentHandlerOption {
	return func(h *AgentHandler) { h.writeCrontab = fn }
}

func WithLogReader(fn agents.LogReadFunc) AgentHandlerOption {
	return func(h *AgentHandler) { h.readLog = fn }
}

func NewAgentHandler(readFn agents.ReadFunc, opts ...AgentHandlerOption) *AgentHandler {
	h := &AgentHandler{
		readCrontab:  readFn,
		writeCrontab: agents.WriteCrontab,
		readLog:      DefaultLogReader,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (ah *AgentHandler) HandleAgents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ah.listAgents(w, r)
	case http.MethodPost:
		ah.createAgent(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ah *AgentHandler) HandleAgentAction(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/agents/")
	if path == "" {
		http.Error(w, "missing agent index", http.StatusBadRequest)
		return
	}

	if strings.HasSuffix(path, "/toggle") {
		ah.toggleAgent(w, r, strings.TrimSuffix(path, "/toggle"))
		return
	}
	if strings.HasSuffix(path, "/log") {
		ah.getLog(w, r, strings.TrimSuffix(path, "/log"))
		return
	}

	switch r.Method {
	case http.MethodDelete:
		ah.deleteAgent(w, r, path)
	case http.MethodPut:
		ah.updateAgent(w, r, path)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ah *AgentHandler) listAgents(w http.ResponseWriter, r *http.Request) {
	raw, err := ah.readCrontab()
	if err != nil {
		http.Error(w, "failed to read crontab", http.StatusInternalServerError)
		return
	}

	ct := agents.ParseCrontab(raw)
	agentList := ct.Agents()

	out := make([]AgentJSON, len(agentList))
	for i, a := range agentList {
		out[i] = agentToJSON(a)
		out[i].LineIndex = i
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AgentsResponse{Agents: out})
}

func (ah *AgentHandler) resolveCrontabAndAgent(indexStr string) (*agents.Crontab, *agents.Agent, int, error) {
	idx, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("invalid agent index")
	}

	raw, err := ah.readCrontab()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to read crontab")
	}

	ct := agents.ParseCrontab(raw)
	agentList := ct.Agents()
	if idx < 0 || idx >= len(agentList) {
		return nil, nil, 0, fmt.Errorf("agent not found")
	}

	return ct, agentList[idx], idx, nil
}

func (ah *AgentHandler) createAgent(w http.ResponseWriter, r *http.Request) {
	var req AgentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Schedule == "" || req.Directory == "" || req.Harness == "" {
		http.Error(w, "schedule, directory, and harness are required", http.StatusBadRequest)
		return
	}

	raw, err := ah.readCrontab()
	if err != nil {
		http.Error(w, "failed to read crontab", http.StatusInternalServerError)
		return
	}

	ct := agents.ParseCrontab(raw)
	newAgent := &agents.Agent{
		Schedule:  req.Schedule,
		Directory: req.Directory,
		Harness:   agents.Harness(req.Harness),
		Model:     req.Model,
		Prompt:    req.Prompt,
		LogPath:   req.LogPath,
		Enabled:   true,
	}
	ct.AddAgent(newAgent)

	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agentToJSON(newAgent))
}

func (ah *AgentHandler) deleteAgent(w http.ResponseWriter, r *http.Request, indexStr string) {
	ct, _, _, err := ah.resolveCrontabAndAgent(indexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	agentList := ct.Agents()
	idx, _ := strconv.Atoi(indexStr)
	ct.DeleteAgent(agentList[idx].LineIndex)

	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"ok":true}`)
}

func (ah *AgentHandler) toggleAgent(w http.ResponseWriter, r *http.Request, indexStr string) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ct, _, idx, err := ah.resolveCrontabAndAgent(indexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	agentList := ct.Agents()
	ct.ToggleAgent(agentList[idx].LineIndex)
	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	agentList = ct.Agents()
	var agent *agents.Agent
	for _, a := range agentList {
		if a.LineIndex == agentList[idx].LineIndex || a == agentList[idx] {
			agent = a
			break
		}
	}
	if agent == nil && len(agentList) > idx {
		agent = agentList[idx]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentToJSON(agent))
}

func (ah *AgentHandler) updateAgent(w http.ResponseWriter, r *http.Request, indexStr string) {
	ct, existing, _, err := ah.resolveCrontabAndAgent(indexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req AgentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	updated := &agents.Agent{
		Schedule:  req.Schedule,
		Directory: req.Directory,
		Harness:   agents.Harness(req.Harness),
		Model:     req.Model,
		Prompt:    req.Prompt,
		LogPath:   req.LogPath,
		Enabled:   existing.Enabled,
	}
	ct.UpdateAgent(existing.LineIndex, updated)

	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentToJSON(updated))
}

func (ah *AgentHandler) getLog(w http.ResponseWriter, r *http.Request, indexStr string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, agent, _, err := ah.resolveCrontabAndAgent(indexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if agent.LogPath == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&agents.LogInfo{Exists: false})
		return
	}

	logInfo, err := ah.readLog(agent.LogPath, 50)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&agents.LogInfo{Exists: false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logInfo)
}

func agentToJSON(a *agents.Agent) AgentJSON {
	return AgentJSON{
		LineIndex: a.LineIndex,
		Schedule:  a.Schedule,
		Directory: a.Directory,
		Harness:   string(a.Harness),
		Model:     a.Model,
		Prompt:    a.Prompt,
		LogPath:   a.LogPath,
		Enabled:   a.Enabled,
	}
}

func DefaultLogReader(path string, n int) (*agents.LogInfo, error) {
	return agents.ReadLogFile(path, n)
}
