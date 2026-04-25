# Implementation Plan: Improved Search & Filtering

## Phase 1: Search Infrastructure

### Tasks

- [x] 1.1: Write unit tests for SearchIndex data structure [commit: search-index-tests]
- [x] 1.2: Implement SearchIndex with commit message and filename indexing [commit: search-index-impl]
- [x] 1.3: Write unit tests for indexer (concurrent-safe updates) [commit: indexer-tests]
- [x] 1.4: Implement incremental index update on repo pull [commit: indexer-incremental]

## Phase 2: Search API Endpoint

### Tasks

- [x] 2.1: Write unit tests for search endpoint with mock index [commit: search-handler-tests]
- [x] 2.2: Implement GET /api/search endpoint with query parsing [commit: search-endpoint-impl]
- [x] 2.3: Add filtering support (repo, date, author) [commit: search-filtering]
- [x] 2.4: Add result ranking and context snippet generation [commit: search-ranking]

## Phase 3: Search UI

### Tasks

- [ ] 3.1: Write unit tests for search debouncing in JS
- [ ] 3.2: Add search input component to main dashboard header
- [ ] 3.3: Implement filter UI (repo selector, date range, author)
- [ ] 3.4: Connect search UI to API endpoint

## Phase 4: Results Display

### Tasks

- [ ] 4.1: Write tests for results panel rendering
- [ ] 4.2: Implement collapsible results panel with match highlighting
- [ ] 4.3: Implement click-to-navigate to commit detail view

## Phase 5: Integration & Verification

### Tasks

- [ ] 5.1: Run full test suite (`go test ./...`) — all tests must pass
- [ ] 5.2: Verify `go build` completes without errors
- [ ] 5.3: Manual verification of search responsiveness
- [ ] 5.4: Update tech-debt.md and lessons-learned.md
- [ ] 5.5: Commit with git note