# Implementation Plan: Advanced Search with Server-Side Indexing

## Phase 1: Indexer
- [ ] Task: Build in-memory search index
  - [ ] Write tests for index construction
  - [ ] Index commit messages, file paths, and diff summaries per repo
  - [ ] Tokenize with simple word splitting and lowercase
- [ ] Task: Incremental index updates
  - [ ] Write tests for index refresh on new commits
  - [ ] Hook into existing pull scheduler to re-index after pulls

## Phase 2: Search API & UI
- [ ] Task: Add /api/search endpoint
  - [ ] Write tests for search ranking
  - [ ] Accept query param, return ranked matches
  - [ ] Rank by exact match > prefix > substring
- [ ] Task: Build search UI
  - [ ] Write tests for search input and results list
  - [ ] Global search input in header
  - [ ] Results page with repo, commit, and snippet highlights

## Phase 3: Performance & Polish
- [ ] Task: Add debouncing and result limits
  - [ ] Write tests for debounce behavior
  - [ ] 300ms debounce on input
  - [ ] Cap results at 50 with "show more" option
- [ ] Task: Manual verification
  - [ ] Verify search finds commits and files across multiple repos
