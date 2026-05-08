# Implementation Plan: Cross-Repository Commit Search & Filtering

## Phase 1: Backend Query Engine

- [x] Task: Design commit search query model
  - [x] Write tests for query parameter parsing (`q`, `author`, `since`, `until`, `repo`)
  - [x] Implement `CommitSearchQuery` struct with validation
  - [x] Implement relative date parser (`1d`, `7d`, `30d`)
- [x] Task: Build commit search repository layer
  - [x] Write tests for date filtering and pagination
  - [x] Implement dynamic WHERE filtering (repo, author, since, until)
  - [x] Add pagination support (limit/offset)
  - [x] Ensure Timestamp field on CommitDoc for date filtering

## Phase 2: API Endpoint

- [x] Task: Implement `GET /api/commits/search`
  - [x] Write handler tests with mocked repo layer
  - [x] Wire query parser to repo search
  - [x] Return JSON with commits, total count, and pagination metadata
  - [x] Add CORS and basic input sanitization

## Phase 3: Frontend Search UI

- [x] Task: Build search input component
  - [x] Write tests for search bar rendering and events
  - [x] Add text input with debounced input handling (300ms)
  - [x] Add author dropdown populated from known authors
- [x] Task: Build date filter controls
  - [x] Write tests for date filter state changes
  - [x] Add preset buttons: Today, Last 7 Days, Last 30 Days
  - [x] Add optional custom date range inputs
- [x] Task: Connect to API and render results
  - [x] Write tests for fetch integration and result rendering
  - [x] Fetch `/api/commits/search` on filter change
  - [x] Render results in existing commit feed cards
  - [x] Show "No results" state and result count badge

## Phase 4: Persistence & Polish

- [x] Task: Persist recent searches
  - [x] Write tests for localStorage serialization
  - [x] Save last 5 searches to localStorage
  - [x] Show recent searches as quick-tap chips
- [x] Task: Integration & verification
  - [x] Run full test suite (`go test ./...`) — all tests must pass
  - [x] Verify `go build` completes without errors
  - [ ] Manual browser verification of filter combinations
  - [x] Update tech-debt.md and lessons-learned.md
  - [x] Commit with git note
