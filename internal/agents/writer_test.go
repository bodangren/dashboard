package agents

import (
	"strings"
	"testing"
)

func TestWriteCrontab(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()
	agents[0].Model = "gpt-4o-mini"

	result := ct.String()

	if !strings.Contains(result, "gpt-4o-mini") {
		t.Error("updated model should appear in output")
	}
	if !strings.Contains(result, "SHELL=/bin/bash") {
		t.Error("env vars should be preserved")
	}
	if !strings.Contains(result, "/usr/bin/cleanup.sh") {
		t.Error("non-agent lines should be preserved")
	}
	if strings.Count(result, "SHELL=/bin/bash") != 1 {
		t.Error("env vars should not be duplicated")
	}
}

func TestWriteRealWorldFormat(t *testing.T) {
	ct := ParseCrontab(realWorldCrontab)
	agents := ct.Agents()
	agents[0].Model = "openai/gpt-5.4"

	result := ct.String()

	if !strings.Contains(result, "-m openai/gpt-5.4") {
		t.Error("should use -m flag for model")
	}
	if !strings.Contains(result, "run conductor/autonomous_prompt.md") {
		t.Error("should use 'run' subcommand for prompt")
	}
	if !strings.Contains(result, "> ") {
		t.Error("should use single > for log redirect")
	}
}

func TestToggleAgent(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()

	if !agents[0].Enabled {
		t.Fatal("first agent should start enabled")
	}
	if agents[1].Enabled {
		t.Fatal("second agent should start disabled")
	}

	ct.ToggleAgent(agents[0].LineIndex)
	if agents[0].Enabled {
		t.Error("first agent should now be disabled")
	}
	line := ct.Lines[agents[0].LineIndex]
	if !strings.HasPrefix(strings.TrimSpace(line.Raw), "#") {
		t.Error("disabled line should be commented out")
	}

	ct.ToggleAgent(agents[1].LineIndex)
	if !agents[1].Enabled {
		t.Error("second agent should be enabled after toggle")
	}
}

func TestToggleAgentFlipsDisabledToEnabled(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()
	disabled := agents[1]

	ct.ToggleAgent(disabled.LineIndex)
	if !disabled.Enabled {
		t.Error("toggling a disabled agent should enable it")
	}
	line := ct.Lines[disabled.LineIndex]
	if strings.HasPrefix(strings.TrimSpace(line.Raw), "#") {
		t.Error("enabled line should not be commented")
	}
	if !strings.Contains(line.Raw, "gemini") {
		t.Error("enabled line should still contain harness command")
	}
}

func TestAddAgent(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	before := len(ct.Agents())

	newAgent := &Agent{
		Schedule:  "0 6 * * *",
		Directory: "/home/user/projects/lib",
		Harness:   HarnessCodex,
		Model:     "codex-1",
		Prompt:    "fix.md",
		LogPath:   "/var/log/agent-lib.log",
		Enabled:   true,
	}
	ct.AddAgent(newAgent)

	after := ct.Agents()
	if len(after) != before+1 {
		t.Fatalf("expected %d agents after add, got %d", before+1, len(after))
	}

	added := after[len(after)-1]
	if added.Harness != HarnessCodex {
		t.Errorf("expected codex harness, got %q", added.Harness)
	}

	result := ct.String()
	if !strings.Contains(result, "codex") {
		t.Error("crontab output should contain the new agent")
	}
	if !strings.Contains(result, "SHELL=/bin/bash") {
		t.Error("existing content should still be present")
	}
	if !strings.Contains(result, "# /home/user/projects/lib") {
		t.Error("section header comment with directory path should be in output")
	}
}

func TestAddAgentWithOpenCodeFormat(t *testing.T) {
	ct := ParseCrontab("SHELL=/bin/bash\n")

	newAgent := &Agent{
		Schedule:   "30 3,7,11,15,23 * * *",
		Directory:  "/home/daniel-bo/Desktop/kanban",
		Harness:    HarnessOpenCode,
		BinaryPath: "/home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode",
		Model:      "zai-coding-plan/glm-5.1",
		Prompt:     "conductor/autonomous_prompt.md",
		LogPath:    "/home/daniel-bo/Desktop/kanban/conductor/log.md",
		Enabled:    true,
	}
	ct.AddAgent(newAgent)

	result := ct.String()
	if !strings.Contains(result, "/home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode") {
		t.Error("should contain full binary path")
	}
	if !strings.Contains(result, "-m zai-coding-plan/glm-5.1") {
		t.Error("should contain -m flag with model")
	}
	if !strings.Contains(result, "run conductor/autonomous_prompt.md") {
		t.Error("should contain 'run' subcommand with prompt path")
	}
	if !strings.Contains(result, "> /home/daniel-bo/Desktop/kanban/conductor/log.md 2>&1") {
		t.Error("should contain single > redirect to log path")
	}
	if !strings.Contains(result, "# /home/daniel-bo/Desktop/kanban") {
		t.Error("should contain section header with directory path")
	}
}

func TestDeleteAgent(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()
	before := len(agents)

	ct.DeleteAgent(agents[0].LineIndex)

	after := ct.Agents()
	if len(after) != before-1 {
		t.Fatalf("expected %d agents after delete, got %d", before-1, len(after))
	}

	result := ct.String()
	if strings.Contains(result, "opencode") {
		t.Error("deleted agent should not appear in output")
	}
	if !strings.Contains(result, "SHELL=/bin/bash") {
		t.Error("non-agent lines should be preserved after delete")
	}
	if !strings.Contains(result, "/usr/bin/cleanup.sh") {
		t.Error("non-agent cron lines should be preserved after delete")
	}
}

func TestDeleteAgentRemovesSectionHeader(t *testing.T) {
	input := `SHELL=/bin/bash

# My Project
0 */4 * * * cd /home/user/proj && opencode -m openai/gpt-5.4 run t.md > /log/a.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()

	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}

	ct.DeleteAgent(agents[0].LineIndex)
	result := ct.String()

	if strings.Contains(result, "My Project") {
		t.Error("section header comment should be removed with agent")
	}
	if !strings.Contains(result, "SHELL=/bin/bash") {
		t.Error("non-agent lines should be preserved")
	}
}

// ── New tests: ReorganizeAutomation & section grouping ─────────

func TestReorganizeAutomation(t *testing.T) {
	input := `SHELL=/bin/bash
PATH=/usr/bin

# legacy comment
0 5 * * * /usr/bin/backup.sh

# ==AUTOMATION BELOW THIS LINE==
# old header
30 */4 * * * cd /home/user/proj-alpha && opencode -m gpt-4o run t.md > /log/a.log 2>&1

0 8 * * * cd /home/user/proj-beta && opencode -m gpt-5 run d.md > /log/b.log 2>&1
`
	ct := ParseCrontab(input)
	projects := []string{"/home/user/proj-alpha", "/home/user/proj-beta", "/home/user/proj-gamma"}
	ct.ReorganizeAutomation(projects)
	result := ct.String()

	// Separator preserved
	if !strings.Contains(result, "==AUTOMATION BELOW THIS LINE==") {
		t.Error("separator should be in output")
	}

	// Preamble preserved
	if !strings.Contains(result, "SHELL=/bin/bash") {
		t.Error("preamble env var should be preserved")
	}
	if !strings.Contains(result, "/usr/bin/backup.sh") {
		t.Error("preamble cron job should be preserved")
	}

	// All projects have section headers
	if !strings.Contains(result, "# /home/user/proj-alpha") {
		t.Error("proj-alpha should have section header")
	}
	if !strings.Contains(result, "# /home/user/proj-beta") {
		t.Error("proj-beta should have section header")
	}
	if !strings.Contains(result, "# /home/user/proj-gamma") {
		t.Error("proj-gamma (empty) should have section header")
	}

	// Agents grouped under correct sections
	alphaIdx := strings.Index(result, "# /home/user/proj-alpha")
	betaIdx := strings.Index(result, "# /home/user/proj-beta")
	gammaIdx := strings.Index(result, "# /home/user/proj-gamma")
	agentAlphaIdx := strings.Index(result, "cd /home/user/proj-alpha")
	agentBetaIdx := strings.Index(result, "cd /home/user/proj-beta")

	if agentAlphaIdx < alphaIdx || agentAlphaIdx > betaIdx {
		t.Error("proj-alpha agent should be between its header and proj-beta header")
	}
	if agentBetaIdx < betaIdx || agentBetaIdx > gammaIdx {
		t.Error("proj-beta agent should be between its header and proj-gamma header")
	}

	// Old header comment is gone
	if strings.Contains(result, "# old header") {
		t.Error("old stray header should be cleaned up")
	}
}

func TestAddAgentUnderExistingSection(t *testing.T) {
	input := `SHELL=/bin/bash

# ==AUTOMATION BELOW THIS LINE==
# /home/user/proj
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md > /log/a.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}

	ct.AddAgent(&Agent{
		Schedule:  "0 8 * * *",
		Directory: "/home/user/proj",
		Harness:   HarnessOpenCode,
		Model:     "gpt-5",
		Prompt:    "d.md",
		LogPath:   "/log/b.log",
		Enabled:   true,
	})

	result := ct.String()

	// Only one section header for the project
	if strings.Count(result, "# /home/user/proj") != 1 {
		t.Errorf("should have exactly one section header, got:\n%s", result)
	}

	// Both agents present
	if !strings.Contains(result, "gpt-4o") {
		t.Error("first agent should be present")
	}
	if !strings.Contains(result, "gpt-5") {
		t.Error("second agent should be present")
	}

	// Both agents appear after the header
	headerIdx := strings.Index(result, "# /home/user/proj")
	agent1Idx := strings.Index(result, "gpt-4o")
	agent2Idx := strings.Index(result, "gpt-5")
	if agent1Idx < headerIdx || agent2Idx < headerIdx {
		t.Error("both agents should appear after section header")
	}
}

func TestDeleteAgentPreservesHeader(t *testing.T) {
	input := `SHELL=/bin/bash

# ==AUTOMATION BELOW THIS LINE==
# /home/user/proj
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md > /log/a.log 2>&1
0 8 * * * cd /home/user/proj && opencode -m gpt-5 run d.md > /log/b.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()
	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}

	// Delete first agent — header should stay because second agent remains
	ct.DeleteAgent(agents[0].LineIndex)

	result := ct.String()
	if !strings.Contains(result, "# /home/user/proj") {
		t.Error("section header should be preserved when other agents remain")
	}
	if strings.Contains(result, "gpt-4o") {
		t.Error("deleted agent should not appear")
	}
	if !strings.Contains(result, "gpt-5") {
		t.Error("remaining agent should still appear")
	}
}

func TestDeleteLastAgentRemovesHeader(t *testing.T) {
	input := `SHELL=/bin/bash

# ==AUTOMATION BELOW THIS LINE==
# /home/user/proj
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md > /log/a.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}

	ct.DeleteAgent(agents[0].LineIndex)

	result := ct.String()
	if strings.Contains(result, "# /home/user/proj") {
		t.Error("section header should be removed when last agent is deleted")
	}
	if strings.Contains(result, "gpt-4o") {
		t.Error("deleted agent should not appear")
	}
}

func TestUpdateAgentMovesProject(t *testing.T) {
	input := `SHELL=/bin/bash

# ==AUTOMATION BELOW THIS LINE==
# /home/user/proj-alpha
0 */4 * * * cd /home/user/proj-alpha && opencode -m gpt-4o run t.md > /log/a.log 2>&1

# /home/user/proj-beta
`
	ct := ParseCrontab(input)
	agents := ct.Agents()
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}

	// Move agent from proj-alpha to proj-beta
	updated := &Agent{
		Schedule:  agents[0].Schedule,
		Directory: "/home/user/proj-beta",
		Harness:   agents[0].Harness,
		Model:     "gpt-5",
		Prompt:    agents[0].Prompt,
		LogPath:   agents[0].LogPath,
		Enabled:   agents[0].Enabled,
	}
	ct.UpdateAgent(agents[0].LineIndex, updated)

	result := ct.String()

	// Agent should now be under proj-beta
	if strings.Contains(result, "cd /home/user/proj-alpha") {
		t.Error("agent should no longer be under proj-alpha")
	}
	if !strings.Contains(result, "cd /home/user/proj-beta") {
		t.Error("agent should now be under proj-beta")
	}
	if !strings.Contains(result, "gpt-5") {
		t.Error("updated model should appear")
	}

	// proj-alpha header should remain (it's in the projects list via ReorganizeAutomation)
	// but proj-beta header should exist
	if !strings.Contains(result, "# /home/user/proj-beta") {
		t.Error("proj-beta section header should exist")
	}
}

func TestReorganizeCreatesSeparator(t *testing.T) {
	ct := ParseCrontab("SHELL=/bin/bash\n")
	ct.ReorganizeAutomation([]string{"/home/user/proj"})

	result := ct.String()
	if !strings.Contains(result, "==AUTOMATION BELOW THIS LINE==") {
		t.Error("separator should be created when missing")
	}
	if !strings.Contains(result, "# /home/user/proj") {
		t.Error("project header should be created")
	}
}

func TestReorganizePreservesContentAboveSeparator(t *testing.T) {
	input := `SHELL=/bin/bash
HOME=/home/user

# maintenance
0 3 * * * /usr/bin/cleanup.sh

# ==AUTOMATION BELOW THIS LINE==
# /home/user/proj
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md > /log/a.log 2>&1
`
	ct := ParseCrontab(input)
	ct.ReorganizeAutomation([]string{"/home/user/proj"})
	result := ct.String()

	// Everything above separator preserved
	if !strings.Contains(result, "SHELL=/bin/bash") {
		t.Error("env vars above separator should be preserved")
	}
	if !strings.Contains(result, "HOME=/home/user") {
		t.Error("HOME env var should be preserved")
	}
	if !strings.Contains(result, "/usr/bin/cleanup.sh") {
		t.Error("maintenance cron should be preserved")
	}

	// Verify order: preamble before separator
	shellIdx := strings.Index(result, "SHELL=")
	sepIdx := strings.Index(result, "AUTOMATION BELOW THIS LINE")
	if shellIdx > sepIdx {
		t.Error("preamble should appear before separator")
	}
}
