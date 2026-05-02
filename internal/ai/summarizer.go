package ai

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

var ErrNoAPIKey = errors.New("OPENAI_API_KEY not set")

type CommitInfo struct {
	Hash    string
	Author  string
	Message string
	Body    string
}

type SummarizerConfig struct {
	APIKey     string
	CacheTTL   time.Duration
	MaxRetries int
}

type CommitSummary struct {
	RepoPath   string
	CommitHash string
	Summary    string
	Flags      []string
	CachedAt   time.Time
}

type LLMProvider interface {
	Summarize(ctx context.Context, prompt string) (string, error)
}

type summarizer struct {
	provider LLMProvider
	cache    map[string]*CommitSummary
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

var (
	defaultSummarizer *summarizer
	defaultMu         sync.Once
)

func newSummarizer(cfg SummarizerConfig) (*summarizer, error) {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 5 * time.Minute
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}

	var provider LLMProvider
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		provider = &mockProvider{}
	} else {
		provider = &openAIProvider{apiKey: apiKey}
	}

	return &summarizer{
		provider: provider,
		cache:    make(map[string]*CommitSummary),
		cacheTTL: cfg.CacheTTL,
	}, nil
}

func DefaultSummarizer() (*summarizer, error) {
	var err error
	defaultMu.Do(func() {
		defaultSummarizer, err = newSummarizer(SummarizerConfig{})
	})
	return defaultSummarizer, err
}

func (s *summarizer) SummarizeCommits(ctx context.Context, repoPath string, commits []CommitInfo) (string, error) {
	if len(commits) == 0 {
		return "", nil
	}

	cacheKey := cacheKey(repoPath, commits[0].Hash)

	s.cacheMu.RLock()
	if cached := s.getFromCache(cacheKey); cached != nil && time.Since(cached.CachedAt) < s.cacheTTL {
		s.cacheMu.RUnlock()
		return cached.Summary, nil
	}
	s.cacheMu.RUnlock()

	prompt := buildPrompt(repoPath, commits)
	summary, err := s.provider.Summarize(ctx, prompt)

	s.cacheMu.Lock()
	s.cache[cacheKey] = &CommitSummary{
		RepoPath:   repoPath,
		CommitHash: commits[0].Hash,
		Summary:    summary,
		Flags:      detectFlags(commits),
		CachedAt:   time.Now(),
	}
	s.cacheMu.Unlock()

	return summary, err
}

func (s *summarizer) getFromCache(key string) *CommitSummary {
	return s.cache[key]
}

func cacheKey(repoPath, commitHash string) string {
	return repoPath + ":" + commitHash
}

func buildPrompt(repoPath string, commits []CommitInfo) string {
	var sb strings.Builder
	sb.WriteString("You are a code reviewer analyzing git commits.\n\n")
	sb.WriteString("Repository: ")
	sb.WriteString(repoPath)
	sb.WriteString("\n\nRecent commits:\n\n")

	for i, c := range commits {
		if i >= 10 {
			break
		}
		sb.WriteString("- ")
		sb.WriteString(c.Hash)
		sb.WriteString(" | ")
		sb.WriteString(c.Author)
		sb.WriteString(" | ")
		sb.WriteString(c.Message)
		sb.WriteString("\n")
	}

	sb.WriteString("\nProvide a brief summary (1-2 sentences) of what changed in these commits.")
	return sb.String()
}

func detectFlags(commits []CommitInfo) []string {
	var flags []string
	for _, c := range commits {
		if strings.Contains(c.Body, "<<<<<<<") || strings.Contains(c.Body, "=======") || strings.Contains(c.Body, ">>>>>>>") {
			flags = append(flags, "conflict")
		}
	}
	if len(commits) > 5 {
		flags = append(flags, "rapid-changes")
	}
	return flags
}

func DetectCommitFlags(commit CommitInfo) []string {
	var flags []string
	if strings.Contains(commit.Body, "<<<<<<<") || strings.Contains(commit.Body, "=======") || strings.Contains(commit.Body, ">>>>>>>") {
		flags = append(flags, "conflict-marker")
	}
	if strings.Contains(commit.Message, "WIP") || strings.Contains(commit.Message, "work in progress") {
		flags = append(flags, "wip")
	}
	return flags
}

type mockProvider struct{}

func (p *mockProvider) Summarize(ctx context.Context, prompt string) (string, error) {
	return "Mock summary: Several commits across the repository.", nil
}

type openAIProvider struct {
	apiKey string
}

func (p *openAIProvider) Summarize(ctx context.Context, prompt string) (string, error) {
	return "", ErrNoAPIKey
}

type ActivityEnhancer struct {
	summarizer *summarizer
	eventCache map[string]*CommitSummary
	cacheMu    sync.RWMutex
}

func NewActivityEnhancer(s *summarizer) *ActivityEnhancer {
	return &ActivityEnhancer{
		summarizer: s,
		eventCache: make(map[string]*CommitSummary),
	}
}

func (ae *ActivityEnhancer) EnhanceEvent(repoPath, commitHash string, commits []CommitInfo) (string, []string, error) {
	cacheKey := cacheKey(repoPath, commitHash)

	ae.cacheMu.RLock()
	if cached := ae.eventCache[cacheKey]; cached != nil {
		ae.cacheMu.RUnlock()
		return cached.Summary, cached.Flags, nil
	}
	ae.cacheMu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	summary, err := ae.summarizer.SummarizeCommits(ctx, repoPath, commits)
	if err != nil {
		return "", nil, err
	}

	ae.cacheMu.Lock()
	ae.eventCache[cacheKey] = &CommitSummary{
		RepoPath:   repoPath,
		CommitHash: commitHash,
		Summary:    summary,
		Flags:      detectFlags(commits),
		CachedAt:   time.Now(),
	}
	ae.cacheMu.Unlock()

	return summary, detectFlags(commits), nil
}