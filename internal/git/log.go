package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Commit represents a single git commit.
type Commit struct {
	Hash      string    // 7-character short hash
	Message   string    // subject line
	Body      string    // commit body (may be empty)
	Notes     string    // git notes (may be empty)
	Author    string    // author name
	Timestamp time.Time // commit timestamp (UTC)
}

// fieldSep separates fields within a commit record in git log output.
const fieldSep = "\x1f"

// recordSep separates commit records in git log output.
// Must be a printable string so it can be passed as a process argument.
const recordSep = "---GIT-RECORD-END---"

// GetCommits returns up to n commits from the git repo at repoPath,
// most recent first.
func GetCommits(repoPath string, n int) ([]Commit, error) {
	// Each commit is: hash FS subject FS body FS author FS unix-ts FS notes recordSep
	// %x1f is the ASCII unit separator — git outputs it as a raw byte.
	// --notes ensures %N is populated from refs/notes/commits.
	format := "%h%x1f%s%x1f%b%x1f%an%x1f%ct%x1f%N" + recordSep

	cmd := exec.Command("git", "log", fmt.Sprintf("-n%d", n), "--notes", "--format="+format)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log in %s: %w", repoPath, err)
	}

	return parseCommits(string(out))
}

// parseCommits splits raw git log output into Commit structs.
func parseCommits(raw string) ([]Commit, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	records := strings.Split(raw, recordSep)
	var commits []Commit
	for _, rec := range records {
		rec = strings.TrimSpace(rec)
		if rec == "" {
			continue
		}
		parts := strings.SplitN(rec, fieldSep, 6)
		if len(parts) < 5 {
			continue
		}
		hash := strings.TrimSpace(parts[0])
		message := strings.TrimSpace(parts[1])
		body := strings.TrimSpace(parts[2])
		author := strings.TrimSpace(parts[3])
		tsRaw := strings.TrimSpace(parts[4])
		notes := ""
		if len(parts) == 6 {
			notes = strings.TrimSpace(parts[5])
		}

		ts, err := parseUnixTimestamp(tsRaw)
		if err != nil {
			continue
		}

		commits = append(commits, Commit{
			Hash:      hash,
			Message:   message,
			Body:      body,
			Notes:     notes,
			Author:    author,
			Timestamp: ts,
		})
	}
	return commits, nil
}

func parseUnixTimestamp(s string) (time.Time, error) {
	unix, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp %q: %w", s, err)
	}
	return time.Unix(unix, 0).UTC(), nil
}
