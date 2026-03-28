package agents

import (
	"testing"
)

const sampleCrontab = `# ┌───────────── minute (0 - 59)
# │ ┌───────────── hour (0 - 23)
SHELL=/bin/bash
PATH=/usr/local/bin:/usr/bin:/bin

# dashboard agent
0 */4 * * * cd /home/user/projects/dashboard && opencode --model gpt-4o --prompt tasks.md >> /var/log/agent-dashboard.log 2>&1

# disabled agent
# 0 8 * * * cd /home/user/projects/api && gemini --model gemini-2.0-flash --prompt daily.md >> /var/log/agent-api.log 2>&1

# unrelated cron
30 2 * * * /usr/bin/cleanup.sh
`

func TestParseCrontab(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}

	a := agents[0]
	if !a.Enabled {
		t.Error("first agent should be enabled")
	}
	if a.Schedule != "0 */4 * * *" {
		t.Errorf("expected schedule '0 */4 * * *', got %q", a.Schedule)
	}
	if a.Directory != "/home/user/projects/dashboard" {
		t.Errorf("expected directory '/home/user/projects/dashboard', got %q", a.Directory)
	}
	if a.Harness != HarnessOpenCode {
		t.Errorf("expected harness opencode, got %q", a.Harness)
	}
	if a.Model != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %q", a.Model)
	}
	if a.LogPath != "/var/log/agent-dashboard.log" {
		t.Errorf("expected log path '/var/log/agent-dashboard.log', got %q", a.LogPath)
	}

	b := agents[1]
	if b.Enabled {
		t.Error("second agent should be disabled (commented out)")
	}
	if b.Schedule != "0 8 * * *" {
		t.Errorf("expected schedule '0 8 * * *', got %q", b.Schedule)
	}
	if b.Harness != HarnessGemini {
		t.Errorf("expected harness gemini, got %q", b.Harness)
	}
	if b.Directory != "/home/user/projects/api" {
		t.Errorf("expected directory '/home/user/projects/api', got %q", b.Directory)
	}
}

func TestParseCrontabPreservesNonAgent(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)

	kinds := make(map[LineKind]int)
	for _, l := range ct.Lines {
		kinds[l.Kind]++
	}

	if kinds[LineComment] < 4 {
		t.Errorf("expected at least 4 comment lines, got %d", kinds[LineComment])
	}
	if kinds[LineEnvVar] != 2 {
		t.Errorf("expected 2 env var lines, got %d", kinds[LineEnvVar])
	}
	if kinds[LineAgent] != 2 {
		t.Errorf("expected 2 agent lines, got %d", kinds[LineAgent])
	}
	if kinds[LineOther] != 1 {
		t.Errorf("expected 1 other line, got %d", kinds[LineOther])
	}
}

func TestParseCrontabCodexAgent(t *testing.T) {
	input := `0 6 * * * cd /home/user/projects/lib && codex --model codex-1 --prompt fix.md >> /var/log/agent-lib.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()

	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].Harness != HarnessCodex {
		t.Errorf("expected harness codex, got %q", agents[0].Harness)
	}
}

func TestParseCrontabEmpty(t *testing.T) {
	ct := ParseCrontab("")
	if len(ct.Lines) != 0 {
		t.Errorf("expected 0 lines for empty input, got %d", len(ct.Lines))
	}
	agents := ct.Agents()
	if len(agents) != 0 {
		t.Errorf("expected 0 agents for empty input, got %d", len(agents))
	}
}
