package ai

import (
	"context"
	"testing"
	"time"
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
	commits := []CommitInfo{
		{Hash: "abc123", Message: "Fix bug in login", Author: "Alice"},
		{Hash: "def456", Message: "Add new feature", Author: "Bob"},
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
	conflictCommits := []CommitInfo{
		{Body: "<<<<<<< HEAD\nsome content\n=======\nother content\n>>>>>>> branch", Hash: "abc", Message: "fix", Author: "a"},
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

	manyCommits := make([]CommitInfo, 6)
	for i := range manyCommits {
		manyCommits[i] = CommitInfo{Hash: "abc", Message: "fix", Author: "a"}
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
		commit   CommitInfo
		wantFlag string
	}{
		{
			name:     "conflict markers in body",
			commit:   CommitInfo{Body: "<<<<<<< HEAD\ncontent\n=======\nother\n>>>>>>> branch"},
			wantFlag: "conflict-marker",
		},
		{
			name:     "WIP in message",
			commit:   CommitInfo{Message: "WIP: implement feature"},
			wantFlag: "wip",
		},
		{
			name:     "work in progress",
			commit:   CommitInfo{Message: "work in progress on auth"},
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

	commits := []CommitInfo{
		{Hash: "abc123", Message: "Initial commit", Author: "Alice"},
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

	commits := []CommitInfo{
		{Hash: "abc123", Message: "Initial commit", Author: "Alice"},
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

	commits := []CommitInfo{
		{Hash: "abc123", Message: "Initial commit", Author: "Alice"},
	}
	summary, flags, err := ae.EnhanceEvent("/test/repo", "abc123", commits)
	if err != nil {
		t.Fatalf("EnhanceEvent() error = %v", err)
	}
	if summary == "" {
		t.Error("summary should not be empty")
	}
	if len(flags) != 0 {
		t.Errorf("flags = %v, want empty", flags)
	}
}

func TestActivityEnhancerCaching(t *testing.T) {
	mp := &testProvider{}
	s := &summarizer{
		provider: mp,
		cache:    make(map[string]*CommitSummary),
		cacheTTL: 5 * time.Minute,
	}
	ae := NewActivityEnhancer(s)

	commits := []CommitInfo{
		{Hash: "abc123", Message: "Initial commit", Author: "Alice"},
	}

	ae.EnhanceEvent("/test/repo", "abc123", commits)
	ae.EnhanceEvent("/test/repo", "abc123", commits)

	if mp.called != 1 {
		t.Errorf("provider called %d times, want 1 (cached)", mp.called)
	}
}

func TestActivityEnhancerFlags(t *testing.T) {
	mp := &testProvider{}
	s := &summarizer{
		provider: mp,
		cache:    make(map[string]*CommitSummary),
		cacheTTL: 5 * time.Minute,
	}
	ae := NewActivityEnhancer(s)

	commits := []CommitInfo{
		{Body: "<<<<<<< HEAD\ncontent\n=======\nother\n>>>>>>> branch", Hash: "abc123", Message: "fix", Author: "Alice"},
	}
	_, flags, err := ae.EnhanceEvent("/test/repo", "abc123", commits)
	if err != nil {
		t.Fatalf("EnhanceEvent() error = %v", err)
	}
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
}