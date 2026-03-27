package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler_callsPullOnEachRepo(t *testing.T) {
	repos := []string{"/repo/a", "/repo/b", "/repo/c"}

	var mu sync.Mutex
	pulled := []string{}
	pullFn := func(path string) error {
		mu.Lock()
		pulled = append(pulled, path)
		mu.Unlock()
		return nil
	}

	s := New(repos, 50*time.Millisecond, 0, pullFn)
	s.RunOnce()

	mu.Lock()
	defer mu.Unlock()
	if len(pulled) != len(repos) {
		t.Fatalf("expected %d pulls, got %d: %v", len(repos), len(pulled), pulled)
	}
	for i, r := range repos {
		if pulled[i] != r {
			t.Errorf("pull[%d]: got %q, want %q", i, pulled[i], r)
		}
	}
}

func TestScheduler_runsSerially(t *testing.T) {
	repos := []string{"/repo/a", "/repo/b"}
	concurrent := 0
	maxConcurrent := 0
	var mu sync.Mutex

	pullFn := func(path string) error {
		mu.Lock()
		concurrent++
		if concurrent > maxConcurrent {
			maxConcurrent = concurrent
		}
		mu.Unlock()

		time.Sleep(10 * time.Millisecond) // simulate work

		mu.Lock()
		concurrent--
		mu.Unlock()
		return nil
	}

	s := New(repos, 50*time.Millisecond, 0, pullFn)
	s.RunOnce()

	if maxConcurrent > 1 {
		t.Errorf("pulls ran concurrently (max concurrent: %d), expected serial", maxConcurrent)
	}
}

func TestScheduler_logsErrorsWithoutStopping(t *testing.T) {
	repos := []string{"/repo/a", "/repo/b"}
	pulled := 0
	pullFn := func(path string) error {
		pulled++
		if path == "/repo/a" {
			return &testError{"simulated pull failure"}
		}
		return nil
	}

	s := New(repos, 50*time.Millisecond, 0, pullFn)
	s.RunOnce()

	if pulled != 2 {
		t.Errorf("expected 2 pull attempts despite error, got %d", pulled)
	}
}

func TestScheduler_tickerFiresAtInterval(t *testing.T) {
	repos := []string{"/repo/a"}
	callCount := 0
	var mu sync.Mutex

	pullFn := func(path string) error {
		mu.Lock()
		callCount++
		mu.Unlock()
		return nil
	}

	s := New(repos, 30*time.Millisecond, 0, pullFn)
	stop := s.Start()
	time.Sleep(95 * time.Millisecond)
	stop()

	mu.Lock()
	defer mu.Unlock()
	// With 30ms interval and 95ms window, expect 3 ticks (at 0, 30, 60, 90ms)
	if callCount < 2 || callCount > 5 {
		t.Errorf("expected ~3 pulls in 95ms at 30ms interval, got %d", callCount)
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
