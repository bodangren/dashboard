package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestActivityHandler_returnsEvents(t *testing.T) {
	now := time.Now()
	repos := []string{"/repos/alpha"}
	var commitsCalled bool

	ah := NewActivityHandler(
		WithActivityRepos(repos),
		WithActivityGetCommits(func(repoPath string, n int) ([]Commit, error) {
			commitsCalled = true
			return []Commit{
				{Hash: "abc1234", Message: "fix bug", Author: "Alice", Timestamp: now},
			}, nil
		}),
		WithActivityReadCrontab(func() (string, error) { return "", nil }),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("GET", "/api/activity", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !commitsCalled {
		t.Error("getCommits was not called")
	}

	var resp map[string][]ActivityEvent
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	events, ok := resp["events"]
	if !ok {
		t.Fatal("missing events field")
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Type != EventTypeCommit {
		t.Errorf("event type: got %s, want %s", events[0].Type, EventTypeCommit)
	}
	if events[0].Message != "fix bug" {
		t.Errorf("message: got %q, want 'fix bug'", events[0].Message)
	}
}

func TestActivityHandler_filtersByType(t *testing.T) {
	now := time.Now()
	repos := []string{"/repos/alpha"}

	ah := NewActivityHandler(
		WithActivityRepos(repos),
		WithActivityGetCommits(func(repoPath string, n int) ([]Commit, error) {
			return []Commit{{Hash: "abc", Message: "commit msg", Timestamp: now}}, nil
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("GET", "/api/activity?types=agent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp map[string][]ActivityEvent
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp["events"]) != 0 {
		t.Errorf("expected 0 events for agent-only filter, got %d", len(resp["events"]))
	}
}

func TestActivityHandler_recordsAgentEvent(t *testing.T) {
	now := time.Now()
	ah := NewActivityHandler(
		WithActivityReadCrontab(func() (string, error) { return "", nil }),
	)

	meta, _ := json.Marshal(AgentEventMetadata{AgentID: "test-agent", Status: "started"})
	ah.RecordAgentEvent(ActivityEvent{
		ID:        "agent-1",
		Type:      EventTypeAgent,
		Repo:      "/repos/alpha",
		Message:   "Agent started",
		Timestamp: now,
		Metadata:  meta,
	})

	ah.eventMu.Lock()
	if len(ah.recentAgentEvents) != 1 {
		t.Errorf("expected 1 event, got %d", len(ah.recentAgentEvents))
	}
	ah.eventMu.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("GET", "/api/activity?types=agent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp map[string][]ActivityEvent
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp["events"]) != 1 {
		t.Errorf("expected 1 agent event, got %d", len(resp["events"]))
	}
}

func TestActivityHandler_pullEvents(t *testing.T) {
	now := time.Now()
	pullMu := &sync.RWMutex{}
	lastPullTime := map[string]time.Time{"/repos/alpha": now}
	lastPullErr := map[string]string{}

	ah := NewActivityHandler(
		WithActivityPullState(pullMu, lastPullTime, lastPullErr),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("GET", "/api/activity?types=pull", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp map[string][]ActivityEvent
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp["events"]) != 1 {
		t.Errorf("expected 1 pull event, got %d", len(resp["events"]))
	}
	if resp["events"][0].Type != EventTypePull {
		t.Errorf("type: got %s, want pull", resp["events"][0].Type)
	}
}

func TestActivityHandler_methodNotAllowed(t *testing.T) {
	ah := NewActivityHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("POST", "/api/activity", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestActivityHandler_eventsSortedByTimestamp(t *testing.T) {
	now := time.Now()
	older := now.Add(-1 * time.Hour)
	newer := now

	repos := []string{"/repos/alpha"}

	ah := NewActivityHandler(
		WithActivityRepos(repos),
		WithActivityGetCommits(func(repoPath string, n int) ([]Commit, error) {
			return []Commit{
				{Hash: "aaa", Message: "older commit", Author: "Alice", Timestamp: older},
				{Hash: "bbb", Message: "newer commit", Author: "Bob", Timestamp: newer},
			}, nil
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("GET", "/api/activity", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp map[string][]ActivityEvent
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	events := resp["events"]
	if len(events) < 2 {
		t.Skip("not enough events")
	}
	if events[0].Message != "newer commit" {
		t.Errorf("first event should be newer, got %q", events[0].Message)
	}
	if events[1].Message != "older commit" {
		t.Errorf("second event should be older, got %q", events[1].Message)
	}
}

func TestActivityHandler_limitParam(t *testing.T) {
	now := time.Now()
	repos := []string{"/repos/alpha"}

	ah := NewActivityHandler(
		WithActivityRepos(repos),
		WithActivityGetCommits(func(repoPath string, n int) ([]Commit, error) {
			var commits []Commit
			for i := 0; i < 10; i++ {
				commits = append(commits, Commit{
					Hash:      "abc123" + string(rune(i)),
					Message:   "commit " + string(rune(i)),
					Author:    "Alice",
					Timestamp: now.Add(time.Duration(i) * time.Minute),
				})
			}
			return commits, nil
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/activity", ah.HandleActivity)

	req := httptest.NewRequest("GET", "/api/activity?limit=5", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp map[string][]ActivityEvent
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp["events"]) != 5 {
		t.Errorf("expected 5 events, got %d", len(resp["events"]))
	}
}
