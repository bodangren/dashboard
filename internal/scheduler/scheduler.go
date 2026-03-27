// Package scheduler runs periodic git pull operations across a set of repos.
package scheduler

import (
	"log"
	"time"
)

// PullFunc is the function called to pull a single repository.
type PullFunc func(repoPath string) error

// Scheduler pulls a list of repos at a fixed interval, serially.
type Scheduler struct {
	repos    []string
	interval time.Duration
	between  time.Duration // wait between repos to avoid hammering remotes
	pullFn   PullFunc
}

// New creates a Scheduler.
//   - repos: list of repo paths to pull
//   - interval: how often to run the full pull cycle (e.g. 4 * time.Hour)
//   - between: sleep between individual repo pulls (0 is fine for tests)
//   - pullFn: function to call for each repo
func New(repos []string, interval, between time.Duration, pullFn PullFunc) *Scheduler {
	return &Scheduler{
		repos:    repos,
		interval: interval,
		between:  between,
		pullFn:   pullFn,
	}
}

// RunOnce pulls every repo once, serially, logging any errors.
func (s *Scheduler) RunOnce() {
	for _, repo := range s.repos {
		if err := s.pullFn(repo); err != nil {
			log.Printf("scheduler: pull failed for %s: %v", repo, err)
		}
		if s.between > 0 {
			time.Sleep(s.between)
		}
	}
}

// Start launches the scheduler in a background goroutine, firing immediately
// and then at each interval. Returns a stop function.
func (s *Scheduler) Start() (stop func()) {
	ticker := time.NewTicker(s.interval)
	done := make(chan struct{})

	go func() {
		s.RunOnce() // fire immediately on start
		for {
			select {
			case <-ticker.C:
				s.RunOnce()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	return func() { close(done) }
}
