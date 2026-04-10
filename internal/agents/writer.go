package agents

import (
	"fmt"
	"os/exec"
	"strings"
)

const automationSeparator = "# ==AUTOMATION BELOW THIS LINE=="

// ── Serialization ──────────────────────────────────────────────

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

func buildAgentLine(a *Agent) string {
	binary := string(a.Harness)
	if a.BinaryPath != "" {
		binary = a.BinaryPath
	}

	var cmd []string
	cmd = append(cmd, binary)
	if a.Model != "" {
		cmd = append(cmd, "-m", a.Model)
	}
	if a.Prompt != "" {
		cmd = append(cmd, "run", a.Prompt)
	}

	line := fmt.Sprintf("%s cd %s && %s", a.Schedule, a.Directory, strings.Join(cmd, " "))

	if a.LogPath != "" {
		line += fmt.Sprintf(" > %s 2>&1", a.LogPath)
	}

	return line
}

// ── Separator helpers ──────────────────────────────────────────

func findSeparator(ct *Crontab) int {
	for i, l := range ct.Lines {
		if l.Kind == LineComment && strings.Contains(l.Raw, "AUTOMATION BELOW THIS LINE") {
			return i
		}
	}
	return -1
}

func (c *Crontab) ensureSeparator() {
	if findSeparator(c) >= 0 {
		return
	}
	c.Lines = append(c.Lines, Line{
		Raw:  automationSeparator,
		Kind: LineComment,
	})
}

// ── ReorganizeAutomation ───────────────────────────────────────

// ReorganizeAutomation rebuilds everything below the automation separator:
// one section header per project directory, with agents grouped underneath.
// Content above the separator is preserved untouched.
func (c *Crontab) ReorganizeAutomation(projects []string) {
	c.ensureSeparator()
	sepIdx := findSeparator(c)

	// Collect all agents currently below the separator (deep copy to avoid aliasing)
	var existingAgents []*Agent
	for i := sepIdx + 1; i < len(c.Lines); i++ {
		if c.Lines[i].Agent != nil {
			a := *c.Lines[i].Agent
			existingAgents = append(existingAgents, &a)
		}
	}

	// Keep only preamble + separator
	c.Lines = append([]Line(nil), c.Lines[:sepIdx+1]...)

	// Group agents by directory
	agentsByDir := make(map[string][]*Agent)
	for _, a := range existingAgents {
		agentsByDir[a.Directory] = append(agentsByDir[a.Directory], a)
	}

	// Track directories already placed (from the projects list)
	placed := make(map[string]bool)

	// Build automation section: project headers first, then orphan dirs
	for _, dir := range projects {
		placed[dir] = true
		c.Lines = append(c.Lines, Line{Raw: "# " + dir, Kind: LineComment})
		for _, a := range agentsByDir[dir] {
			a.LineIndex = len(c.Lines)
			c.Lines = append(c.Lines, Line{
				Raw:   agentRawLine(a),
				Kind:  LineAgent,
				Agent: a,
			})
		}
	}

	// Agents whose directory isn't in the projects list
	for dir, agents := range agentsByDir {
		if placed[dir] {
			continue
		}
		c.Lines = append(c.Lines, Line{Raw: "# " + dir, Kind: LineComment})
		for _, a := range agents {
			a.LineIndex = len(c.Lines)
			c.Lines = append(c.Lines, Line{
				Raw:   agentRawLine(a),
				Kind:  LineAgent,
				Agent: a,
			})
		}
	}
}

// ── Mutations ──────────────────────────────────────────────────

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
	c.ensureSeparator()
	// Append agent to Lines so ReorganizeAutomation collects it,
	// then rebuild the automation section with proper grouping.
	c.Lines = append(c.Lines, Line{
		Raw:   agentRawLine(a),
		Kind:  LineAgent,
		Agent: a,
	})
	c.ReorganizeAutomation([]string{})
}

func (c *Crontab) DeleteAgent(lineIndex int) {
	if lineIndex < 0 || lineIndex >= len(c.Lines) {
		return
	}
	if c.Lines[lineIndex].Agent == nil {
		return
	}

	// Check if preceding line is a section header comment
	commentIdx := -1
	if lineIndex > 0 && c.Lines[lineIndex-1].Kind == LineComment {
		commentIdx = lineIndex - 1
	}

	// Check if next line is another agent in the same section
	// (same directory with no comment between them)
	nextIsSameSection := false
	if lineIndex+1 < len(c.Lines) && commentIdx >= 0 {
		next := c.Lines[lineIndex+1]
		if next.Agent != nil && next.Agent.Directory == c.Lines[lineIndex].Agent.Directory {
			nextIsSameSection = true
		}
	}

	startIdx := lineIndex
	if commentIdx >= 0 && !nextIsSameSection {
		startIdx = commentIdx
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

	oldDir := line.Agent.Directory

	if updated.Directory == oldDir {
		// Same project — update in place
		updated.LineIndex = lineIndex
		line.Agent = updated
		line.Raw = agentRawLine(updated)
	} else {
		// Different project — remove from old section, append to new
		c.DeleteAgent(lineIndex)
		c.Lines = append(c.Lines, Line{
			Raw:   agentRawLine(updated),
			Kind:  LineAgent,
			Agent: updated,
		})
		c.ReorganizeAutomation(nil)
	}
}

// ── System crontab I/O ────────────────────────────────────────

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

// ── Helpers ────────────────────────────────────────────────────

func agentRawLine(a *Agent) string {
	if a.Enabled {
		return buildAgentLine(a)
	}
	return "# " + buildAgentLine(a)
}

func isAgentLineKind(l Line) bool {
	return l.Agent != nil
}
