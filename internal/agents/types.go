package agents

import (
	"fmt"
	"time"
)

type ReadFunc func() (string, error)
type WriteFunc func(string) error
type LogReadFunc func(path string, n int) (*LogInfo, error)

type Harness string

const (
	HarnessOpenCode Harness = "opencode"
	HarnessGemini   Harness = "gemini"
	HarnessCodex    Harness = "codex"
)

type Agent struct {
	Schedule      string  // cron expression (e.g. "0 */4 * * *")
	Directory     string  // working directory for the agent
	Harness       Harness // opencode, gemini, or codex
	BinaryPath    string  // full path to binary (e.g. /home/user/.nvm/.../opencode)
	Model         string  // model name (e.g. "zai-coding-plan/glm-5.1")
	Prompt        string  // prompt file path (positional after `run`)
	LogPath       string  // path to log file
	SectionHeader string  // comment line above this agent (section label)
	Enabled       bool    // false if the line is commented out
	LineIndex     int     // position in the Crontab.Lines slice
}

func (a *Agent) AgentID() string {
	return fmt.Sprintf("%s:%s:%s", a.Schedule, a.Directory, a.Model)
}

type LineKind int

const (
	LineBlank LineKind = iota
	LineComment
	LineEnvVar
	LineAgent
	LineOther
)

type Line struct {
	Raw     string
	Kind    LineKind
	Agent   *Agent
	Comment string
}

type Crontab struct {
	Lines []Line
}

func (c *Crontab) Agents() []*Agent {
	var out []*Agent
	for i := range c.Lines {
		if c.Lines[i].Agent != nil {
			out = append(out, c.Lines[i].Agent)
		}
	}
	return out
}

func (c *Crontab) AgentByIndex(lineIndex int) *Agent {
	if lineIndex >= 0 && lineIndex < len(c.Lines) {
		return c.Lines[lineIndex].Agent
	}
	return nil
}

func (c *Crontab) AgentByID(id string) *Agent {
	for _, a := range c.Agents() {
		if a.AgentID() == id {
			return a
		}
	}
	return nil
}

type LogInfo struct {
	Exists    bool      `json:"exists"`
	LastRun   time.Time `json:"last_run"`
	Lines     []string  `json:"lines"`
	Truncated bool      `json:"truncated"`
}

type AgentState struct {
	ExitCode  int
	LastError string
}

type AgentStateMap struct {
	states map[string]*AgentState
}

func NewAgentStateMap() *AgentStateMap {
	return &AgentStateMap{states: make(map[string]*AgentState)}
}

func (m *AgentStateMap) Get(agentID string) *AgentState {
	if s, ok := m.states[agentID]; ok {
		return s
	}
	return nil
}

func (m *AgentStateMap) Set(agentID string, state *AgentState) {
	m.states[agentID] = state
}

func (m *AgentStateMap) Clear(agentID string) {
	delete(m.states, agentID)
}
