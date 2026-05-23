# Implementation Plan: Commit Bookmarking & Personal Annotations

## Phase 1: Data Model & Storage (TDD)

- [ ] Task: Define `Bookmark` type and SQLite schema.
  - [ ] Write tests for `Bookmark` serialization/deserialization.
  - [ ] Create migration for `bookmarks` table.
- [ ] Task: Implement `BookmarkStore` with Create, Get, List, UpdateNote, Delete.
  - [ ] Write unit tests for all CRUD methods.
  - [ ] Write tests for list filtering by repo and search by note content.
  - [ ] Implement store with `database/sql`.
- [ ] Task: Add `BookmarkStore` to `HandlerConfig`.
  - [ ] Write tests confirming handlers receive store via config.
- [ ] Task: Measure — User Manual Verification 'Phase 1: Data Model & Storage'

## Phase 2: API Endpoints (TDD)

- [ ] Task: Implement `POST /api/bookmarks` — create bookmark.
  - [ ] Write handler tests for success, duplicate hash (idempotent), missing params.
- [ ] Task: Implement `GET /api/bookmarks` — list bookmarks with `?repo=` and `?q=` filters.
  - [ ] Write handler tests for filtering, search, and empty results.
- [ ] Task: Implement `PUT /api/bookmarks/:id` — update note.
  - [ ] Write handler tests for update and not-found cases.
- [ ] Task: Implement `DELETE /api/bookmarks/:id` — remove bookmark.
  - [ ] Write handler tests for delete and not-found cases.
- [ ] Task: Measure — User Manual Verification 'Phase 2: API Endpoints'

## Phase 3: Frontend — Bookmark Toggle & Note Editor (TDD)

- [ ] Task: Add bookmark star toggle to commit cards in commit feed.
  - [ ] Write JS tests for toggle click handler and API call.
  - [ ] Implement toggle with visual state (filled vs outline star).
- [ ] Task: Add inline note editor to commit cards (expand on bookmark).
  - [ ] Write JS tests for note save/cancel flow.
  - [ ] Implement textarea with save/cancel buttons.
- [ ] Task: Wire bookmark state into diff view (`diff.html`).
  - [ ] Manual verification: bookmark and annotate from diff page.
- [ ] Task: Measure — User Manual Verification 'Phase 3: Frontend Bookmark Toggle'

## Phase 4: Bookmarks Page (TDD)

- [ ] Task: Create `/bookmarks` HTML page and route.
  - [ ] Fetch `/api/bookmarks` on load.
  - [ ] Render bookmark list with repo name, commit message, timestamp, and note.
- [ ] Task: Add repo filter dropdown and search input.
  - [ ] Write JS tests for filter and search behavior.
  - [ ] Implement client-side filtering with debounced search.
- [ ] Task: Add "Go to diff" link for each bookmarked commit.
- [ ] Task: Add bookmark count badge to project cards on main dashboard.
  - [ ] Write JS tests for badge rendering.
- [ ] Task: Measure — User Manual Verification 'Phase 4: Bookmarks Page'

## Phase 5: Integration & Polish

- [ ] Task: Ensure bookmarked commits are visually distinguished in activity feed.
- [ ] Task: Add CSS for star icon, note editor, and bookmarks page layout (mobile-first).
- [ ] Task: Run full Go test suite with `-race` flag.
- [ ] Task: Update `tracks.md` and commit.
- [ ] Task: Measure — User Manual Verification 'Phase 5: Integration & Polish'
