package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dashboard/internal/agents"
	"dashboard/internal/ws"
)

type AgentJSON struct {
	ID            string `json:"id"`
	LineIndex     int    `json:"line_index"`
	Schedule      string `json:"schedule"`
	Directory     string `json:"directory"`
	Harness       string `json:"harness"`
	BinaryPath    string `json:"binary_path"`
	Model         string `json:"model"`
	Prompt        string `json:"prompt"`
	LogPath       string `json:"log_path"`
	SectionHeader string `json:"section_header"`
	Enabled       bool   `json:"enabled"`
	LastError     string `json:"last_error,omitempty"`
	ExitCode      int    `json:"exit_code,omitempty"`
}

type AgentsResponse struct {
	Agents []AgentJSON `json:"agents"`
}

type AgentCreateRequest struct {
	Schedule      string `json:"schedule"`
	Directory     string `json:"directory"`
	Harness       string `json:"harness"`
	BinaryPath    string `json:"binary_path"`
	Model         string `json:"model"`
	Prompt        string `json:"prompt"`
	LogPath       string `json:"log_path"`
	SectionHeader string `json:"section_header"`
}

type AgentHandler struct {
	readCrontab  agents.ReadFunc
	writeCrontab agents.WriteFunc
	readLog      agents.LogReadFunc
	repos        []string
	openCodeBin  string
	watcherMgr   *ws.WatcherManager
	stateMap     *agents.AgentStateMap
}

type AgentHandlerOption func(*AgentHandler)

func WithWriteFunc(fn agents.WriteFunc) AgentHandlerOption {
	return func(h *AgentHandler) { h.writeCrontab = fn }
}

func WithLogReader(fn agents.LogReadFunc) AgentHandlerOption {
	return func(h *AgentHandler) { h.readLog = fn }
}

func WithOpenCodeBinary(path string) AgentHandlerOption {
	return func(h *AgentHandler) { h.openCodeBin = path }
}

func WithWatcherManager(wm *ws.WatcherManager) AgentHandlerOption {
	return func(h *AgentHandler) { h.watcherMgr = wm }
}

func WithAgentStateMap(sm *agents.AgentStateMap) AgentHandlerOption {
	return func(h *AgentHandler) { h.stateMap = sm }
}

func NewAgentHandler(readFn agents.ReadFunc, opts ...AgentHandlerOption) *AgentHandler {
	h := &AgentHandler{
		readCrontab:  readFn,
		writeCrontab: agents.WriteCrontab,
		readLog:      DefaultLogReader,
		stateMap:     agents.NewAgentStateMap(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// SetRepos configures the project directories used for crontab section headers.
func (ah *AgentHandler) SetRepos(repos []string) {
	ah.repos = repos
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
	if strings.HasSuffix(path, "/trigger") {
		ah.triggerAgent(w, r, strings.TrimSuffix(path, "/trigger"))
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
		out[i] = ah.agentToJSON(a)
		out[i].LineIndex = i
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AgentsResponse{Agents: out})
}

func (ah *AgentHandler) resolveCrontabAndAgent(id string) (*agents.Crontab, *agents.Agent, error) {
	raw, err := ah.readCrontab()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read crontab")
	}

	ct := agents.ParseCrontab(raw)
	agent := ct.AgentByID(id)
	if agent == nil {
		return nil, nil, fmt.Errorf("agent not found")
	}

	return ct, agent, nil
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
		Schedule:      req.Schedule,
		Directory:     req.Directory,
		Harness:       agents.Harness(req.Harness),
		BinaryPath:    req.BinaryPath,
		Model:         req.Model,
		Prompt:        req.Prompt,
		LogPath:       req.LogPath,
		SectionHeader: req.SectionHeader,
		Enabled:       true,
	}
	ct.AddAgent(newAgent)
	ct.ReorganizeAutomation(ah.repos)

	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ah.agentToJSON(newAgent))
}

func (ah *AgentHandler) deleteAgent(w http.ResponseWriter, r *http.Request, id string) {
	ct, agent, err := ah.resolveCrontabAndAgent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ct.DeleteAgent(agent.LineIndex)
	ct.ReorganizeAutomation(ah.repos)

	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"ok":true}`)
}

func (ah *AgentHandler) toggleAgent(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ct, agent, err := ah.resolveCrontabAndAgent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ct.ToggleAgent(agent.LineIndex)
	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ah.agentToJSON(agent))
}

func (ah *AgentHandler) triggerAgent(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, agent, err := ah.resolveCrontabAndAgent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	go ah.runAgentAsync(agent, ah.stateMap)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "triggered", "agent_id": id})
}

func (ah *AgentHandler) runAgentAsync(a *agents.Agent, stateMap *agents.AgentStateMap) {
	if ah.watcherMgr != nil && a.LogPath != "" {
		ah.watcherMgr.StartWatching(a.AgentID(), a.LogPath)
	}

	binary := string(a.Harness)
	if a.BinaryPath != "" {
		binary = a.BinaryPath
	}

	var args []string
	if a.Model != "" {
		args = append(args, "-m", a.Model)
	}
	if a.Prompt != "" {
		args = append(args, "run", a.Prompt)
	}

	cmd := exec.Command(binary, args...)
	cmd.Dir = a.Directory

	stderrBuf := &bytes.Buffer{}
	cmd.Stderr = stderrBuf

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	if exitCode != 0 && stateMap != nil {
		stateMap.Set(a.AgentID(), &agents.AgentState{
			ExitCode:  exitCode,
			LastError: stderrBuf.String(),
		})
	} else if stateMap != nil && exitCode == 0 {
		stateMap.Clear(a.AgentID())
	}

	if ah.watcherMgr != nil {
		ah.watcherMgr.StopWatching(a.AgentID())
	}
}

func (ah *AgentHandler) updateAgent(w http.ResponseWriter, r *http.Request, id string) {
	ct, existing, err := ah.resolveCrontabAndAgent(id)
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
		Schedule:      req.Schedule,
		Directory:     req.Directory,
		Harness:       agents.Harness(req.Harness),
		BinaryPath:    req.BinaryPath,
		Model:         req.Model,
		Prompt:        req.Prompt,
		LogPath:       req.LogPath,
		SectionHeader: req.SectionHeader,
		Enabled:       existing.Enabled,
	}
	ct.UpdateAgent(existing.LineIndex, updated)
	ct.ReorganizeAutomation(ah.repos)

	if err := ah.writeCrontab(ct.String()); err != nil {
		http.Error(w, "failed to write crontab", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ah.agentToJSON(updated))
}

func (ah *AgentHandler) getLog(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, agent, err := ah.resolveCrontabAndAgent(id)
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

func (ah *AgentHandler) agentToJSON(a *agents.Agent) AgentJSON {
	aj := AgentJSON{
		ID:            a.AgentID(),
		LineIndex:     a.LineIndex,
		Schedule:      a.Schedule,
		Directory:     a.Directory,
		Harness:       string(a.Harness),
		BinaryPath:    a.BinaryPath,
		Model:         a.Model,
		Prompt:        a.Prompt,
		LogPath:       a.LogPath,
		SectionHeader: a.SectionHeader,
		Enabled:       a.Enabled,
	}
	if ah.stateMap != nil {
		if state := ah.stateMap.Get(a.AgentID()); state != nil {
			aj.ExitCode = state.ExitCode
			aj.LastError = state.LastError
		}
	}
	return aj
}

func DefaultLogReader(path string, n int) (*agents.LogInfo, error) {
	return agents.ReadLogFile(path, n)
}

// discoverModels runs `opencode models` and returns the available model list.
func discoverModels(binaryPath string) []string {
	bin := binaryPath
	if bin == "" {
		bin = resolveOpenCodeBinary()
	}
	if bin == "" {
		return nil
	}
	out, err := exec.Command(bin, "models").Output()
	if err != nil {
		return nil
	}
	return parseModelsOutput(string(out))
}

func parseModelsOutput(raw string) []string {
	var models []string
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			models = append(models, line)
		}
	}
	return models
}

func resolveOpenCodeBinary() string {
	if path, err := exec.LookPath("opencode"); err == nil {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	glob := filepath.Join(home, ".nvm", "versions", "node", "*", "bin", "opencode")
	matches, _ := filepath.Glob(glob)
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return ""
}

func (ah *AgentHandler) HandleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	models := discoverModels(ah.openCodeBin)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"models": models})
}
