package metrics

import (
	"database/sql"
	"os"
	"testing"
)

func TestRecordRun(t *testing.T) {
	db := mustOpenDB(t)
	defer os.Remove(db.Path())

	run := &AgentRun{
		AgentID:     "0 */4 * * *:/home/user/proj:gpt-4o",
		ProjectPath: "/home/user/proj",
		StartedAt:   parseTime("2026-05-09T10:00:00Z"),
		DurationMs:  15000,
		Success:     true,
		ErrorMsg:    "",
	}

	err := db.RecordRun(run)
	if err != nil {
		t.Fatalf("RecordRun failed: %v", err)
	}

	runs, err := db.QueryRuns("0 */4 * * *:/home/user/proj:gpt-4o", 0, "")
	if err != nil {
		t.Fatalf("QueryRuns failed: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].DurationMs != 15000 {
		t.Errorf("DurationMs: got %d, want 15000", runs[0].DurationMs)
	}
	if !runs[0].Success {
		t.Error("Success: got false, want true")
	}
}

func TestRecordRun_Failed(t *testing.T) {
	db := mustOpenDB(t)
	defer os.Remove(db.Path())

	run := &AgentRun{
		AgentID:     "0 */4 * * *:/home/user/proj:gpt-4o",
		ProjectPath: "/home/user/proj",
		StartedAt:   parseTime("2026-05-09T10:00:00Z"),
		DurationMs:  5000,
		Success:     false,
		ErrorMsg:    "command not found",
	}

	err := db.RecordRun(run)
	if err != nil {
		t.Fatalf("RecordRun failed: %v", err)
	}

	runs, err := db.QueryRuns("0 */4 * * *:/home/user/proj:gpt-4o", 0, "")
	if err != nil {
		t.Fatalf("QueryRuns failed: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].Success {
		t.Error("Success: got true, want false")
	}
	if runs[0].ErrorMsg != "command not found" {
		t.Errorf("ErrorMsg: got %q, want 'command not found'", runs[0].ErrorMsg)
	}
}

func TestQueryRuns_Limit(t *testing.T) {
	db := mustOpenDB(t)
	defer os.Remove(db.Path())

	for i := 0; i < 5; i++ {
		run := &AgentRun{
			AgentID:     "test-agent",
			ProjectPath: "/home/proj",
			StartedAt:   parseTime("2026-05-09T10:00:00Z").Add(int64(i) * 3600_000_000),
			DurationMs:  int64(i * 1000),
			Success:     i%2 == 0,
		}
		if err := db.RecordRun(run); err != nil {
			t.Fatalf("RecordRun %d failed: %v", i, err)
		}
	}

	runs, err := db.QueryRuns("test-agent", 3, "")
	if err != nil {
		t.Fatalf("QueryRuns failed: %v", err)
	}
	if len(runs) != 3 {
		t.Errorf("expected 3 runs, got %d", len(runs))
	}
}

func TestQueryRuns_ProjectFilter(t *testing.T) {
	db := mustOpenDB(t)
	defer os.Remove(db.Path())

	run1 := &AgentRun{
		AgentID:     "test-agent",
		ProjectPath: "/home/proj1",
		StartedAt:   parseTime("2026-05-09T10:00:00Z"),
		DurationMs:  1000,
		Success:     true,
	}
	run2 := &AgentRun{
		AgentID:     "test-agent",
		ProjectPath: "/home/proj2",
		StartedAt:   parseTime("2026-05-09T11:00:00Z"),
		DurationMs:  2000,
		Success:     true,
	}
	db.RecordRun(run1)
	db.RecordRun(run2)

	runs, err := db.QueryRuns("test-agent", 0, "/home/proj1")
	if err != nil {
		t.Fatalf("QueryRuns failed: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].ProjectPath != "/home/proj1" {
		t.Errorf("ProjectPath: got %q, want '/home/proj1'", runs[0].ProjectPath)
	}
}

func TestQueryRuns_Empty(t *testing.T) {
	db := mustOpenDB(t)
	defer os.Remove(db.Path())

	runs, err := db.QueryRuns("nonexistent-agent", 0, "")
	if err != nil {
		t.Fatalf("QueryRuns failed: %v", err)
	}
	if len(runs) != 0 {
		t.Errorf("expected 0 runs, got %d", len(runs))
	}
}

func mustOpenDB(t *testing.T) *DB {
	t.Helper()
	f, err := os.CreateTemp("", "metrics_test_*.db")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	f.Close()
	db, err := OpenDB(f.Name())
	if err != nil {
		os.Remove(f.Name())
		t.Fatalf("OpenDB failed: %v", err)
	}
	return db
}