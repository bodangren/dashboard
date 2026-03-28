package agents

import "time"

type Harness string

const (
	HarnessOpenCode Harness = "opencode"
	HarnessGemini   Harness = "gemini"
	HarnessCodex    Harness = "codex"
)

type Agent struct {
	Schedule  string  // cron expression (e.g. "0 */4 * * *")
	Directory string  // working directory for the agent
	Harness   Harness // opencode, gemini, or codex
	Model     string  // model name (e.g. "gpt-4o", "gemini-2.0-flash")
	Prompt    string  // prompt file path or inline command
	LogPath   string  // path to log file
	Enabled   bool    // false if the line is commented out
	LineIndex int     // position in the Crontab.Lines slice
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

type LogInfo struct {
	Exists    bool      `json:"exists"`
	LastRun   time.Time `json:"last_run"`
	Lines     []string  `json:"lines"`
	Truncated bool      `json:"truncated"`
}
