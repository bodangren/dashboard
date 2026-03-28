package agents

import (
	"regexp"
	"strings"
)

var harnessPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bopencode\b`),
	regexp.MustCompile(`\bgemini\b`),
	regexp.MustCompile(`\bcodex\b`),
}

var cronLineRe = regexp.MustCompile(
	`^#?\s*` +
		`(\S+\s+\S+\s+\S+\s+\S+\s+\S+)` + // cron schedule (5 fields)
		`\s+` +
		`(.*)$`, // rest of the line
)

var cdRe = regexp.MustCompile(`cd\s+(\S+)`)
var modelRe = regexp.MustCompile(`--model\s+(\S+)`)
var promptRe = regexp.MustCompile(`--prompt\s+(\S+)`)
var logRedirectRe = regexp.MustCompile(`>>\s*(\S+)\s*2>&1`)

func ParseCrontab(raw string) *Crontab {
	if raw == "" {
		return &Crontab{}
	}

	lines := strings.Split(raw, "\n")
	ct := &Crontab{Lines: make([]Line, 0, len(lines))}

	for i, raw := range lines {
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineBlank})
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			if isCommentedAgentLine(trimmed) {
				agent := extractAgent(trimmed, i, false)
				ct.Lines = append(ct.Lines, Line{
					Raw:   raw,
					Kind:  LineAgent,
					Agent: agent,
				})
				continue
			}
			ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineComment, Comment: trimmed})
			continue
		}

		if isEnvVarLine(trimmed) {
			ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineEnvVar})
			continue
		}

		if isAgentLine(trimmed) {
			agent := extractAgent(trimmed, i, true)
			ct.Lines = append(ct.Lines, Line{
				Raw:   raw,
				Kind:  LineAgent,
				Agent: agent,
			})
			continue
		}

		ct.Lines = append(ct.Lines, Line{Raw: raw, Kind: LineOther})
	}

	return ct
}

func isAgentLine(line string) bool {
	for _, p := range harnessPatterns {
		if p.MatchString(line) {
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
	return strings.Contains(line, "=") && !strings.Contains(line, " ") || regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`).MatchString(line)
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
	model := extractFirstMatch(modelRe, rest)
	prompt := extractFirstMatch(promptRe, rest)
	logPath := extractFirstMatch(logRedirectRe, rest)

	return &Agent{
		Schedule:  schedule,
		Directory: dir,
		Harness:   harness,
		Model:     model,
		Prompt:    prompt,
		LogPath:   logPath,
		Enabled:   enabled,
		LineIndex: lineIndex,
	}
}

func detectHarness(line string) Harness {
	for _, p := range harnessPatterns {
		if p.MatchString(line) {
			return Harness(p.String()[2 : len(p.String())-2])
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
