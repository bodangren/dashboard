package ws

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogWatcher_BasicCheck(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	hub := NewHub()
	hub.Start()
	defer hub.Stop()

	watcher := NewLogWatcher(hub, "agent-1", logPath)

	if _, err := os.Create(logPath); err != nil {
		t.Fatal(err)
	}

	watcher.check()

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != watcher.lastSize {
		t.Errorf("expected lastSize %d, got %d", info.Size(), watcher.lastSize)
	}
}

func TestLogWatcher_InferLineType(t *testing.T) {
	tests := []struct {
		line  string
		wants string
	}{
		{"Everything normal", "stdout"},
		{"ERROR: something failed", "stderr"},
		{"FATAL: crash", "stderr"},
		{"warning: low memory", "warn"},
		{"[WARN] disk full", "warn"},
		{"panic: out of bounds", "stderr"},
		{"INFO: started", "stdout"},
	}

	for _, tc := range tests {
		got := inferLineType(tc.line)
		if got != tc.wants {
			t.Errorf("inferLineType(%q) = %q, want %q", tc.line, got, tc.wants)
		}
	}
}

func TestWatcherManager_StartStop(t *testing.T) {
	hub := NewHub()
	hub.Start()
	defer hub.Stop()

	wm := NewWatcherManager(hub)

	wm.StartWatching("agent-1", "/nonexistent/path.log")

	time.Sleep(10 * time.Millisecond)

	wm.StopWatching("agent-1")
}

func TestWatcherManager_StopAll(t *testing.T) {
	hub := NewHub()
	hub.Start()
	defer hub.Stop()

	wm := NewWatcherManager(hub)

	wm.StartWatching("agent-1", "/nonexistent/path1.log")
	wm.StartWatching("agent-2", "/nonexistent/path2.log")

	time.Sleep(10 * time.Millisecond)

	wm.StopAll()
}
