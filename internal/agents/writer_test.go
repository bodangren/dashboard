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
