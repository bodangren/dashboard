package search

import (
	"strings"
	"sync"
	"time"
	"unicode"
)

type CommitDoc struct {
	RepoPath  string
	Hash      string
	Message   string
	Author    string
	Files     []string
	Timestamp time.Time
}

type SearchResult struct {
	CommitDoc
	Score float64
}

type SearchIndex struct {
	mu     sync.RWMutex
	docs   map[string]CommitDoc
	tokens map[string][]string
}

func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		docs:   make(map[string]CommitDoc),
		tokens: make(map[string][]string),
	}
}

func tokenize(text string) []string {
	text = strings.ToLower(text)
	var tokens []string
	var current strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

func (idx *SearchIndex) AddCommit(doc CommitDoc) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	key := doc.RepoPath + ":" + doc.Hash
	idx.docs[key] = doc

	tokens := tokenize(doc.Message)
	for _, f := range doc.Files {
		tokens = append(tokens, tokenize(f)...)
	}
	for _, t := range tokens {
		idx.tokens[t] = append(idx.tokens[t], key)
	}
}

func (idx *SearchIndex) Search(query string) []SearchResult {
	if query == "" {
		return nil
	}

	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tokens := tokenize(query)
	scoreMap := make(map[string]float64)

	for _, tok := range tokens {
		for key := range idx.docs {
			doc := idx.docs[key]
			score := 0.0
			msgTokens := tokenize(doc.Message)
			for _, mt := range msgTokens {
				if mt == tok {
					score += 2.0
				} else if strings.Contains(mt, tok) {
					score += 1.0
				}
			}
			for _, f := range doc.Files {
				if strings.Contains(strings.ToLower(f), tok) {
					score += 1.5
				}
			}
			if score > 0 {
				scoreMap[key] += score
			}
		}
	}

	var results []SearchResult
	for key, score := range scoreMap {
		results = append(results, SearchResult{CommitDoc: idx.docs[key], Score: score})
	}

	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results
}

func (idx *SearchIndex) SearchWithFilters(query, repoPath, author, dateFrom string) []SearchResult {
	return idx.SearchWithQuery(&CommitSearchQuery{
		Q:      query,
		Repo:   repoPath,
		Author: author,
	})
}

func (idx *SearchIndex) SearchWithQuery(q *CommitSearchQuery) []SearchResult {
	if err := q.Validate(); err != nil {
		return nil
	}
	q.SetDefaults()

	results := idx.Search(q.Q)

	var filtered []SearchResult
	for _, r := range results {
		if q.Repo != "" && r.RepoPath != q.Repo {
			continue
		}
		if q.Author != "" && r.Author != q.Author {
			continue
		}
		if q.Since != nil && r.Timestamp.Before(*q.Since) {
			continue
		}
		if q.Until != nil && r.Timestamp.After(*q.Until) {
			continue
		}
		filtered = append(filtered, r)
	}

	total := len(filtered)
	if q.Offset >= total {
		return []SearchResult{}
	}

	end := q.Offset + q.Limit
	if end > total {
		end = total
	}

	return filtered[q.Offset:end]
}
