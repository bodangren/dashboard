package agents

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadLogFile_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	content := "line1\nline2\nline3\n"
	if err := os.WriteFile(logPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}

	info, err := ReadLogFile(logPath, 10)
	if err != nil {
		t.Fatalf("ReadLogFile returned error: %v", err)
	}
	if !info.Exists {
		t.Error("expected Exists=true for existing file")
	}
	if len(info.Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(info.Lines))
	}
	if info.Truncated {
		t.Error("expected Truncated=false for small file")
	}
}

func TestReadLogFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "empty.log")
	if err := os.WriteFile(logPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}

	info, err := ReadLogFile(logPath, 10)
	if err != nil {
		t.Fatalf("ReadLogFile returned error: %v", err)
	}
	if !info.Exists {
		t.Error("expected Exists=true for empty file")
	}
	if len(info.Lines) != 0 {
		t.Errorf("expected 0 lines for empty file, got %d", len(info.Lines))
	}
}

func TestReadLogFile_FileNotFound(t *testing.T) {
	info, err := ReadLogFile("/nonexistent/path/log.log", 10)
	if err != nil {
		t.Fatalf("ReadLogFile should not return error for non-existent file, got: %v", err)
	}
	if info.Exists {
		t.Error("expected Exists=false for non-existent file")
	}
}

func TestReadLogFile_Truncation(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "truncated.log")
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = "line"
	}
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(logPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}

	info, err := ReadLogFile(logPath, 10)
	if err != nil {
		t.Fatalf("ReadLogFile returned error: %v", err)
	}
	if !info.Truncated {
		t.Error("expected Truncated=true for large file with small n")
	}
	if len(info.Lines) != 10 {
		t.Errorf("expected 10 lines after truncation, got %d", len(info.Lines))
	}
}

func TestReadLogFile_ModTimeSet(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "mtime.log")
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)
	if err := os.WriteFile(logPath, []byte("content\n"), 0644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}
	os.Chtimes(logPath, pastTime, pastTime)

	info, err := ReadLogFile(logPath, 10)
	if err != nil {
		t.Fatalf("ReadLogFile returned error: %v", err)
	}
	if info.LastRun.IsZero() {
		t.Error("expected LastRun to be set")
	}
}

func TestAgentByIndex(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()

	if len(agents) == 0 {
		t.Fatal("expected agents from sampleCrontab")
	}

	lineIdx := agents[0].LineIndex
	found := ct.AgentByIndex(lineIdx)
	if found == nil {
		t.Error("AgentByIndex should return agent for valid index")
	}
	if found != agents[0] {
		t.Error("AgentByIndex should return correct agent")
	}

	nilResult := ct.AgentByIndex(-1)
	if nilResult != nil {
		t.Error("AgentByIndex should return nil for negative index")
	}

	nilResult = ct.AgentByIndex(len(ct.Lines))
	if nilResult != nil {
		t.Error("AgentByIndex should return nil for out-of-bounds index")
	}
}

func TestAgentByID(t *testing.T) {
	ct := ParseCrontab(sampleCrontab)
	agents := ct.Agents()

	if len(agents) == 0 {
		t.Fatal("expected agents from sampleCrontab")
	}

	id := agents[0].AgentID()
	found := ct.AgentByID(id)
	if found == nil {
		t.Error("AgentByID should return agent for valid ID")
	}
	if found != agents[0] {
		t.Error("AgentByID should return correct agent")
	}

	nilResult := ct.AgentByID("nonexistent:id:model")
	if nilResult != nil {
		t.Error("AgentByID should return nil for non-existent ID")
	}
}

func TestAgents_ReturnsAllEnabledAndDisabled(t *testing.T) {
	input := `SHELL=/bin/bash
# enabled agent
0 */4 * * * cd /home/user/proj && opencode -m gpt-4o run t.md > /log/a.log 2>&1
# disabled agent
# 0 8 * * * cd /home/user/other && opencode -m gpt-5 run d.md > /log/b.log 2>&1
`
	ct := ParseCrontab(input)
	agents := ct.Agents()

	if len(agents) != 2 {
		t.Errorf("expected 2 agents (enabled + disabled), got %d", len(agents))
	}

	enabledCount := 0
	disabledCount := 0
	for _, a := range agents {
		if a.Enabled {
			enabledCount++
		} else {
			disabledCount++
		}
	}
	if enabledCount != 1 {
		t.Errorf("expected 1 enabled agent, got %d", enabledCount)
	}
	if disabledCount != 1 {
		t.Errorf("expected 1 disabled agent, got %d", disabledCount)
	}
}
