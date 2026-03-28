package agents

import (
	"fmt"
	"os/exec"
	"strings"
)

func (c *Crontab) String() string {
	var lines []string
	for _, l := range c.Lines {
		if l.Agent != nil {
			if l.Agent.Enabled {
				lines = append(lines, buildAgentLine(l.Agent))
			} else {
				lines = append(lines, "# "+buildAgentLine(l.Agent))
			}
		} else {
			lines = append(lines, l.Raw)
		}
	}
	return strings.Join(lines, "\n")
}

func (c *Crontab) ToggleAgent(lineIndex int) {
	if lineIndex < 0 || lineIndex >= len(c.Lines) {
		return
	}
	line := &c.Lines[lineIndex]
	if line.Agent == nil {
		return
	}

	line.Agent.Enabled = !line.Agent.Enabled

	if line.Agent.Enabled {
		line.Raw = buildAgentLine(line.Agent)
	} else {
		line.Raw = "# " + buildAgentLine(line.Agent)
	}
}

func (c *Crontab) AddAgent(a *Agent) {
	if a.SectionHeader != "" {
		commentLine := Line{
			Raw:     "# " + a.SectionHeader,
			Kind:    LineComment,
			Comment: "# " + a.SectionHeader,
		}
		c.Lines = append(c.Lines, commentLine)
	}

	raw := buildAgentLine(a)
	a.LineIndex = len(c.Lines)
	c.Lines = append(c.Lines, Line{
		Raw:   raw,
		Kind:  LineAgent,
		Agent: a,
	})
}

func (c *Crontab) DeleteAgent(lineIndex int) {
	if lineIndex < 0 || lineIndex >= len(c.Lines) {
		return
	}
	if c.Lines[lineIndex].Agent == nil {
		return
	}

	startIdx := lineIndex
	if startIdx > 0 && c.Lines[startIdx-1].Kind == LineComment {
		startIdx--
	}

	c.Lines = append(c.Lines[:startIdx], c.Lines[lineIndex+1:]...)

	for i := startIdx; i < len(c.Lines); i++ {
		if c.Lines[i].Agent != nil {
			c.Lines[i].Agent.LineIndex = i
		}
	}
}

func (c *Crontab) UpdateAgent(lineIndex int, updated *Agent) {
	if lineIndex < 0 || lineIndex >= len(c.Lines) {
		return
	}
	line := &c.Lines[lineIndex]
	if line.Agent == nil {
		return
	}

	updated.LineIndex = lineIndex
	line.Agent = updated

	if updated.Enabled {
		line.Raw = buildAgentLine(updated)
	} else {
		line.Raw = "# " + buildAgentLine(updated)
	}
}

func buildAgentLine(a *Agent) string {
	var parts []string
	parts = append(parts, a.Schedule)
	parts = append(parts, fmt.Sprintf("cd %s", a.Directory))

	harness := string(a.Harness)
	if a.BinaryPath != "" {
		harness = a.BinaryPath
	}
	parts = append(parts, harness)

	if a.Model != "" {
		parts = append(parts, fmt.Sprintf("-m %s", a.Model))
	}
	if a.Prompt != "" {
		parts = append(parts, fmt.Sprintf("run %s", a.Prompt))
	}
	if a.LogPath != "" {
		parts = append(parts, fmt.Sprintf("> %s 2>&1", a.LogPath))
	}

	return strings.Join(parts, " && ") + " || true"
}

func ReadCrontab() (string, error) {
	cmd := exec.Command("crontab", "-l")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("crontab -l: %w", err)
	}
	return string(out), nil
}

func WriteCrontab(content string) error {
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(content)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("crontab write: %w\n%s", err, string(out))
	}
	return nil
}
