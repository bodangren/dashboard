package agents

import (
	"regexp"
	"strings"
)

type harnessPattern struct {
	re   *regexp.Regexp
	name Harness
}

var harnessPatterns = []harnessPattern{
	{regexp.MustCompile(`\bopencode\b`), HarnessOpenCode},
	{regexp.MustCompile(`\bgemini\b`), HarnessGemini},
	{regexp.MustCompile(`\bcodex\b`), HarnessCodex},
}

var cronLineRe = regexp.MustCompile(
	`^#?\s*` +
		`(\S+\s+\S+\s+\S+\s+\S+\s+\S+)` + // cron schedule (5 fields)
		`\s+` +
		`(.*)$`, // rest of the line
)

var cdRe = regexp.MustCompile(`cd\s+(\S+)`)
var binaryPathRe = regexp.MustCompile(`(/\S*/(?:opencode|gemini|codex)(?:\s|$|\z))`)
var modelShortRe = regexp.MustCompile(`-m\s+(\S+)`)
var modelLongRe = regexp.MustCompile(`--model\s+(\S+)`)
var runPromptRe = regexp.MustCompile(`\brun\s+(\S+)`)
var promptLongRe = regexp.MustCompile(`--prompt\s+(\S+)`)
var logRedirectRe = regexp.MustCompile(`>{1,2}\s*(\S+)\s*2>&1`)
var envVarRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=[^ ]*$`)

func ParseCrontab(raw string) *Crontab {
	if raw == "" {
		return &Crontab{}
	}

	lines := strings.Split(raw, "\n")
	ct := &Crontab{Lines: make([]Line, 0, len(lines))}

	var pendingComment string

	for i, raw := range lines {
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			pendingComment = ""
			ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineBlank})
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			if isCommentedAgentLine(trimmed) {
				agent := extractAgent(trimmed, i, false)
				agent.SectionHeader = pendingComment
				pendingComment = ""
				ct.Lines = append(ct.Lines, Line{
					Raw:   raw,
					Kind:  LineAgent,
					Agent: agent,
				})
				continue
			}
			pendingComment = strings.TrimPrefix(trimmed, "#")
			pendingComment = strings.TrimSpace(pendingComment)
			ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineComment, Comment: trimmed})
			continue
		}

		if isEnvVarLine(trimmed) {
			pendingComment = ""
			ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineEnvVar})
			continue
		}

		if isAgentLine(trimmed) {
			agent := extractAgent(trimmed, i, true)
			agent.SectionHeader = pendingComment
			pendingComment = ""
			ct.Lines = append(ct.Lines, Line{
				Raw:   raw,
				Kind:  LineAgent,
				Agent: agent,
			})
			continue
		}

		pendingComment = ""
		ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineOther})
	}

	return ct
}

func isAgentLine(line string) bool {
	for _, p := range harnessPatterns {
		if p.re.MatchString(line) {
			return true
		}
	}
	return false
}

func isCommentedAgentLine(line string) bool {
	stripped := strings.TrimPrefix(line, "#")
	stripped = strings.TrimLeft(stripped, " \t")
	return isAgentLine(stripped)
}

func isEnvVarLine(line string) bool {
	return envVarRe.MatchString(line)
}

func extractAgent(raw string, lineIndex int, enabled bool) *Agent {
	cleaned := raw
	if !enabled {
		cleaned = strings.TrimPrefix(raw, "#")
		cleaned = strings.TrimLeft(cleaned, " \t")
	}

	m := cronLineRe.FindStringSubmatch(cleaned)
	if m == nil {
		return &Agent{Enabled: enabled, LineIndex: lineIndex}
	}

	schedule := m[1]
	rest := m[2]

	harness := detectHarness(rest)
	dir := extractFirstMatch(cdRe, rest)
	binaryPath := strings.TrimSpace(extractFirstMatch(binaryPathRe, rest))
	model := extractFirstMatch(modelShortRe, rest)
	if model == "" {
		model = extractFirstMatch(modelLongRe, rest)
	}
	prompt := extractFirstMatch(runPromptRe, rest)
	if prompt == "" {
		prompt = extractFirstMatch(promptLongRe, rest)
	}
	logPath := extractFirstMatch(logRedirectRe, rest)

	return &Agent{
		Schedule:   schedule,
		Directory:  dir,
		Harness:    harness,
		BinaryPath: binaryPath,
		Model:      model,
		Prompt:     prompt,
		LogPath:    logPath,
		Enabled:    enabled,
		LineIndex:  lineIndex,
	}
}

func detectHarness(line string) Harness {
	for _, p := range harnessPatterns {
		if p.re.MatchString(line) {
			return p.name
		}
	}
	return ""
}

func extractFirstMatch(re *regexp.Regexp, s string) string {
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}
