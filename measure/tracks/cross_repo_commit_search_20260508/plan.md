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

- [ ] Task: Implement `GET /api/commits/search`
  - [ ] Write handler tests with mocked repo layer
  - [ ] Wire query parser to repo search
  - [ ] Return JSON with commits, total count, and pagination metadata
  - [ ] Add CORS and basic input sanitization

## Phase 3: Frontend Search UI

- [ ] Task: Build search input component
  - [ ] Write tests for search bar rendering and events
  - [ ] Add text input with debounced input handling (300ms)
  - [ ] Add author dropdown populated from known authors
- [ ] Task: Build date filter controls
  - [ ] Write tests for date filter state changes
  - [ ] Add preset buttons: Today, Last 7 Days, Last 30 Days
  - [ ] Add optional custom date range inputs
- [ ] Task: Connect to API and render results
  - [ ] Write tests for fetch integration and result rendering
  - [ ] Fetch `/api/commits/search` on filter change
  - [ ] Render results in existing commit feed cards
  - [ ] Show "No results" state and result count badge

## Phase 4: Persistence & Polish

- [ ] Task: Persist recent searches
  - [ ] Write tests for localStorage serialization
  - [ ] Save last 5 searches to localStorage
  - [ ] Show recent searches as quick-tap chips
- [ ] Task: Integration & verification
  - [ ] Run full test suite (`go test ./...`) — all tests must pass
  - [ ] Verify `go build` completes without errors
  - [ ] Manual browser verification of filter combinations
  - [ ] Update tech-debt.md and lessons-learned.md
  - [ ] Commit with git note
