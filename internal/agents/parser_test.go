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

const realWorldCrontab = `SHELL=/bin/bash
HOME=/home/daniel-bo

# Kanban with Z.ai GLM-5.1
30 3,7,11,15,23 * * * cd /home/daniel-bo/Desktop/kanban-conductor && /home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode -m zai-coding-plan/glm-5.1 run conductor/autonomous_prompt.md > /home/daniel-bo/Desktop/mediarr/conductor/opencode-last-run.log 2>&1

# Advantage Games with OpenAI
45 2,7,12,17,22 * * * cd /home/daniel-bo/Desktop/advantage-games && /home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode -m openai/gpt-5.4-mini run @conductor/autonomous_prompt.md > /home/daniel-bo/Desktop/advantage-games/conductor/opencode-cron.log 2>&1
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

func TestParseRealWorldOpenCodeFormat(t *testing.T) {
	ct := ParseCrontab(realWorldCrontab)
	agents := ct.Agents()

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}

	a := agents[0]
	if a.Harness != HarnessOpenCode {
		t.Errorf("expected harness opencode, got %q", a.Harness)
	}
	if a.Model != "zai-coding-plan/glm-5.1" {
		t.Errorf("expected model 'zai-coding-plan/glm-5.1', got %q", a.Model)
	}
	if a.Prompt != "conductor/autonomous_prompt.md" {
		t.Errorf("expected prompt 'conductor/autonomous_prompt.md', got %q", a.Prompt)
	}
	if a.Directory != "/home/daniel-bo/Desktop/kanban-conductor" {
		t.Errorf("expected directory '/home/daniel-bo/Desktop/kanban-conductor', got %q", a.Directory)
	}
	if a.LogPath != "/home/daniel-bo/Desktop/mediarr/conductor/opencode-last-run.log" {
		t.Errorf("expected log path '/home/daniel-bo/Desktop/mediarr/conductor/opencode-last-run.log', got %q", a.LogPath)
	}
	if a.BinaryPath == "" {
		t.Error("expected non-empty binary path")
	}
	if !a.Enabled {
		t.Error("first agent should be enabled")
	}
	if a.SectionHeader != "Kanban with Z.ai GLM-5.1" {
		t.Errorf("expected section header 'Kanban with Z.ai GLM-5.1', got %q", a.SectionHeader)
	}

	b := agents[1]
	if b.Model != "openai/gpt-5.4-mini" {
		t.Errorf("expected model 'openai/gpt-5.4-mini', got %q", b.Model)
	}
	if b.Prompt != "@conductor/autonomous_prompt.md" {
		t.Errorf("expected prompt '@conductor/autonomous_prompt.md', got %q", b.Prompt)
	}
	if b.SectionHeader != "Advantage Games with OpenAI" {
		t.Errorf("expected section header 'Advantage Games with OpenAI', got %q", b.SectionHeader)
	}
}

func TestParseSectionHeaders(t *testing.T) {
	input := `SHELL=/bin/bash

# My Project Alpha
0 */4 * * * cd /home/user/alpha && opencode -m openai/gpt-5.4-mini run tasks.md > /log/a.log 2>&1

# My Project Beta
0 8 * * * cd /home/user/beta && opencode -m zai-coding-plan/glm-5.1 run daily.md > /log/b.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}
	if agents[0].SectionHeader != "My Project Alpha" {
		t.Errorf("expected section header 'My Project Alpha', got %q", agents[0].SectionHeader)
	}
	if agents[1].SectionHeader != "My Project Beta" {
		t.Errorf("expected section header 'My Project Beta', got %q", agents[1].SectionHeader)
	}
}

func TestParseSingleRedirect(t *testing.T) {
	input := `0 */4 * * * cd /home/user/proj && opencode -m openai/gpt-5.4-mini run tasks.md > /tmp/log.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()

	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].LogPath != "/tmp/log.log" {
		t.Errorf("expected log path '/tmp/log.log', got %q", agents[0].LogPath)
	}
}

func TestParseBinaryPath(t *testing.T) {
	input := `0 */4 * * * cd /home/user/proj && /home/user/.nvm/versions/node/v24.4.0/bin/opencode -m openai/gpt-5.4 run t.md > /log/a.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()

	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].BinaryPath != "/home/user/.nvm/versions/node/v24.4.0/bin/opencode" {
		t.Errorf("expected full binary path, got %q", agents[0].BinaryPath)
	}
}

func TestParseLegacyFormatStillWorks(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}
	if agents[0].Model != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %q", agents[0].Model)
	}
	if agents[0].Prompt != "tasks.md" {
		t.Errorf("expected prompt 'tasks.md', got %q", agents[0].Prompt)
	}
	if agents[0].SectionHeader != "dashboard agent" {
		t.Errorf("expected section header 'dashboard agent', got %q", agents[0].SectionHeader)
	}
}

func TestDetectHarnessExplicitMap(t *testing.T) {
	tests := []struct {
		line string
		want Harness
	}{
		{"opencode -m gpt-4o run tasks.md", HarnessOpenCode},
		{"gemini -m gemini-2.0 run daily.md", HarnessGemini},
		{"codex -m codex-1 run fix.md", HarnessCodex},
		{"opencode run tasks.md", HarnessOpenCode},
		{"some random command with opencode embedded", HarnessOpenCode},
		{"cd /foo && opencode -m test", HarnessOpenCode},
		{"no harness here", ""},
	}

	for _, tc := range tests {
		got := detectHarness(tc.line)
		if got != tc.want {
			t.Errorf("detectHarness(%q) = %q, want %q", tc.line, got, tc.want)
		}
	}
}

func TestIsEnvVarLine_DoesNotMatchCronWithEquals(t *testing.T) {
	tests := []string{
		"0 */4 * * * cd /foo && make test=1",
		"30 2 * * * /usr/bin/cleanup.sh key=value",
		"PATH=/usr/bin:/bin cd /foo && opencode",
	}

	for _, tc := range tests {
		if isEnvVarLine(tc) {
			t.Errorf("isEnvVarLine(%q) = true, want false", tc)
		}
	}
}

func TestIsEnvVarLine_MatchesRealEnvVar(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"SHELL=/bin/bash", true},
		{"PATH=/usr/bin:/bin", true},
		{"_VAR=value", true},
		{"VAR123=value", true},
		{"HOME=/home/user", true},
	}

	for _, tc := range tests {
		got := isEnvVarLine(tc.line)
		if got != tc.want {
			t.Errorf("isEnvVarLine(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}
