package search

import (
	"testing"
)

func TestSearchIndex_AddCommit(t *testing.T) {
	idx := NewSearchIndex()

	commit := CommitDoc{
		RepoPath: "/tmp/testrepo",
		Hash:     "abc1234",
		Message:  "fix authentication bug",
		Author:   "Alice",
		Files:    []string{"auth.go", "login.go"},
	}

	idx.AddCommit(commit)

	results := idx.Search("authentication")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'authentication', got %d", len(results))
	}
	if results[0].Hash != "abc1234" {
		t.Errorf("expected hash abc1234, got %s", results[0].Hash)
	}
}

func TestSearchIndex_SearchByMessage(t *testing.T) {
	idx := NewSearchIndex()

	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "1111111", Message: "add new feature"})
	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "2222222", Message: "fix bug in auth"})
	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo2", Hash: "3333333", Message: "update docs"})

	results := idx.Search("bug")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'bug', got %d", len(results))
	}
	if results[0].Hash != "2222222" {
		t.Errorf("expected hash 2222222, got %s", results[0].Hash)
	}
}

func TestSearchIndex_SearchByFileName(t *testing.T) {
	idx := NewSearchIndex()

	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "1111111", Message: "add feature", Files: []string{"feature.go"}})
	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "2222222", Message: "fix bug", Files: []string{"bugfix.go"}})

	results := idx.Search("feature")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'feature', got %d", len(results))
	}
	if results[0].Hash != "1111111" {
		t.Errorf("expected hash 1111111, got %s", results[0].Hash)
	}
}

func TestSearchIndex_SearchFiltersByRepo(t *testing.T) {
	idx := NewSearchIndex()

	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "1111111", Message: "add feature"})
	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo2", Hash: "2222222", Message: "add feature"})

	results := idx.SearchWithFilters("feature", "/tmp/repo1", "", "")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'feature' in repo1, got %d", len(results))
	}
	if results[0].RepoPath != "/tmp/repo1" {
		t.Errorf("expected repo1, got %s", results[0].RepoPath)
	}
}

func TestSearchIndex_SearchFiltersByAuthor(t *testing.T) {
	idx := NewSearchIndex()

	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "1111111", Message: "fix auth", Author: "Alice"})
	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "2222222", Message: "fix bug", Author: "Bob"})

	author := "Alice"
	results := idx.SearchWithFilters("fix", "", author, "")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'fix' by Alice, got %d", len(results))
	}
	if results[0].Author != "Alice" {
		t.Errorf("expected author Alice, got %s", results[0].Author)
	}
}

func TestSearchIndex_ConcurrentAccess(t *testing.T) {
	idx := NewSearchIndex()

	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			idx.AddCommit(CommitDoc{
				RepoPath: "/tmp/repo1",
				Hash:     "1111111",
				Message:  "fix bug",
			})
		}
		done <- struct{}{}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = idx.Search("fix")
			_ = idx.Search("bug")
		}
		done <- struct{}{}
	}()

	<-done
	<-done
}

func TestSearchIndex_Ranking(t *testing.T) {
	idx := NewSearchIndex()

	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "1111111", Message: "fix critical bug in auth"})
	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "2222222", Message: "fix bug"})

	results := idx.Search("fix bug")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Score < results[1].Score {
		t.Error("results should be ranked by score descending")
	}
}

func TestSearchIndex_EmptyQuery(t *testing.T) {
	idx := NewSearchIndex()

	idx.AddCommit(CommitDoc{RepoPath: "/tmp/repo1", Hash: "1111111", Message: "fix bug"})

	results := idx.Search("")
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}