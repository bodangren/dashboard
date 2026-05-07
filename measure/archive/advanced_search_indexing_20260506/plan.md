# Implementation Plan: Advanced Search with Server-Side Indexing

## Phase 1: Indexer
- [x] Task: Build in-memory search index
  - [x] Write tests for index construction
  - [x] Index commit messages, file paths, and diff summaries per repo
  - [x] Tokenize with simple word splitting and lowercase
- [x] Task: Incremental index updates
  - [x] Write tests for index refresh on new commits
  - [x] Hook into existing pull scheduler to re-index after pulls

## Phase 2: Search API & UI
- [x] Task: Add /api/search endpoint
  - [x] Write tests for search ranking
  - [x] Accept query param, return ranked matches
  - [x] Rank by exact match > prefix > substring
- [x] Task: Build search UI
  - [x] Write tests for search input and results list
  - [x] Global search input in header (already exists)
  - [x] Results page with repo, commit, and snippet highlights
- [x] Task: Wire search indexer into main.go and API handler
  - [x] Indexer builds from repos on startup
  - [x] Background goroutine updates index every 30 minutes
  - [x] searchFunc adapter maps search results to api.SearchResult

## Phase 3: Performance & Polish
- [x] Task: Add debouncing and result limits
  - [x] Write tests for debounce behavior
  - [x] 300ms debounce on input (scheduleSearch)
  - [x] Cap results at 50 with "show more" option
- [ ] Task: Manual verification (deferred - requires browser)
  - [ ] Verify search finds commits and files across multiple repos
