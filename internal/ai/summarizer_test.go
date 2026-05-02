package ai

import (
	"context"
	"testing"
	"time"

	"dashboard/internal/git"
)

type testProvider struct {
	summaries map[string]string
	called    int
}

func (m *testProvider) Summarize(ctx context.Context, prompt string) (string, error) {
	m.called++
	if s, ok := m.summaries[prompt]; ok {
		return s, nil
	}
	return "default mock summary", nil
}

func TestBuildPrompt(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc123", Message: "Fix bug in login", Author: "Alice", Timestamp: time.Now()},
		{Hash: "def456", Message: "Add new feature", Author: "Bob", Timestamp: time.Now()},
	}
	prompt := buildPrompt("/repo/path", commits)

	if prompt == "" {
		t.Fatal("prompt should not be empty")
	}
	if !contains(prompt, "/repo/path") {
		t.Error("prompt should contain repo path")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestDetectFlags(t *testing.T) {
	conflictCommits := []git.Commit{
		{Body: "<<<<<<< HEAD\nsome content\n=======\nother content\n>>>>>>> branch", Hash: "abc", Message: "fix", Author: "a", Timestamp: time.Now()},
	}
	flags := detectFlags(conflictCommits)
	found := false
	for _, f := range flags {
		if f == "conflict" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected conflict flag, got %v", flags)
	}

	manyCommits := make([]git.Commit, 6)
	for i := range manyCommits {
		manyCommits[i] = git.Commit{Hash: "abc", Message: "fix", Author: "a", Timestamp: time.Now()}
	}
	flags = detectFlags(manyCommits)
	found = false
	for _, f := range flags {
		if f == "rapid-changes" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected rapid-changes flag, got %v", flags)
	}
}

func TestDetectCommitFlags(t *testing.T) {
	tests := []struct {
		name     string
		commit   git.Commit
		wantFlag string
	}{
		{
			name:     "conflict markers in body",
			commit:   git.Commit{Body: "<<<<<<< HEAD\ncontent\n=======\nother\n>>>>>>> branch"},
			wantFlag: "conflict-marker",
		},
		{
			name:     "WIP in message",
			commit:   git.Commit{Message: "WIP: implement feature"},
			wantFlag: "wip",
		},
		{
			name:     "work in progress",
			commit:   git.Commit{Message: "work in progress on auth"},
			wantFlag: "wip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := DetectCommitFlags(tt.commit)
			found := false
			for _, f := range flags {
				if f == tt.wantFlag {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected flag %q, got %v", tt.wantFlag, flags)
			}
		})
	}
}

func TestCacheKey(t *testing.T) {
	key := cacheKey("/repo/path", "abc123")
	if key != "/repo/path:abc123" {
		t.Errorf("cacheKey = %q, want %q", key, "/repo/path:abc123")
	}
}

func TestNewSummarizer(t *testing.T) {
	_, err := newSummarizer(SummarizerConfig{})
	if err != nil {
		t.Fatalf("newSummarizer() error = %v", err)
	}
}

func TestSummarizerSummarizeCommits(t *testing.T) {
	mp := &testProvider{
		summaries: make(map[string]string),
	}
	s := &summarizer{
		provider: mp,
		cache:    make(map[string]*CommitSummary),
		cacheTTL: 5 * time.Minute,
	}

	commits := []git.Commit{
		{Hash: "abc123", Message: "Initial commit", Author: "Alice", Timestamp: time.Now()},
	}

	ctx := context.Background()
	summary, err := s.SummarizeCommits(ctx, "/test/repo", commits)
	if err != nil {
		t.Fatalf("SummarizeCommits() error = %v", err)
	}
	if summary == "" {
		t.Error("summary should not be empty")
	}
	if mp.called != 1 {
		t.Errorf("provider called %d times, want 1", mp.called)
	}
}

func TestSummarizerCaching(t *testing.T) {
	mp := &testProvider{}
	s := &summarizer{
		provider: mp,
		cache:    make(map[string]*CommitSummary),
		cacheTTL: 5 * time.Minute,
	}

	commits := []git.Commit{
		{Hash: "abc123", Message: "Initial commit", Author: "Alice", Timestamp: time.Now()},
	}

	ctx := context.Background()

	s.SummarizeCommits(ctx, "/test/repo", commits)
	s.SummarizeCommits(ctx, "/test/repo", commits)

	if mp.called != 1 {
		t.Errorf("provider called %d times, want 1 (cached)", mp.called)
	}
}

func TestActivityEnhancer(t *testing.T) {
	mp := &testProvider{}
	s := &summarizer{
		provider: mp,
		cache:    make(map[string]*CommitSummary),
		cacheTTL: 5 * time.Minute,
	}
	ae := NewActivityEnhancer(s)
	if ae == nil {
		t.Fatal("NewActivityEnhancer returned nil")
	}
}
