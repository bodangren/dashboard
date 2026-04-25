package search

import (
	"dashboard/internal/git"
	"sync"
)

type Indexer struct {
	index       *SearchIndex
	mu          sync.RWMutex
	knownHashes map[string]string
}

func NewIndexer() *Indexer {
	return &Indexer{
		index:       NewSearchIndex(),
		knownHashes: make(map[string]string),
	}
}

func (idx *Indexer) BuildFromRepos(repoPaths []string) error {
	for _, path := range repoPaths {
		if err := idx.indexRepo(path); err != nil {
			return err
		}
	}
	return nil
}

func (idx *Indexer) indexRepo(repoPath string) error {
	commits, err := git.GetCommits(repoPath, 100)
	if err != nil {
		return err
	}

	idx.mu.Lock()
	lastHash, exists := idx.knownHashes[repoPath]
	idx.mu.Unlock()

	if !exists {
		for _, c := range commits {
			idx.index.AddCommit(CommitDoc{
				RepoPath: repoPath,
				Hash:     c.Hash,
				Message:  c.Message,
				Author:   c.Author,
				Files:    nil,
			})
		}
	} else {
		for _, c := range commits {
			if c.Hash == lastHash {
				break
			}
			idx.index.AddCommit(CommitDoc{
				RepoPath: repoPath,
				Hash:     c.Hash,
				Message:  c.Message,
				Author:   c.Author,
				Files:    nil,
			})
		}
	}

	if len(commits) > 0 {
		idx.mu.Lock()
		idx.knownHashes[repoPath] = commits[0].Hash
		idx.mu.Unlock()
	}

	return nil
}

func (idx *Indexer) UpdateRepo(repoPath string) error {
	return idx.indexRepo(repoPath)
}

func (idx *Indexer) Search(query string) []SearchResult {
	return idx.index.Search(query)
}

func (idx *Indexer) SearchWithFilters(query, repoPath, author, dateFrom string) []SearchResult {
	return idx.index.SearchWithFilters(query, repoPath, author, dateFrom)
}