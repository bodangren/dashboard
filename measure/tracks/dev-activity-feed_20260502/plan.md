# Implementation Plan: Unified Developer Activity Feed

## Phase 1: Backend Event Aggregation

### Tasks

- [x] 1.1: Add `Event` type to `internal/api/events.go`
  - Fields: ID, Type (commit|agent|pull), Repo, Message, Timestamp, Metadata (JSON)
- [x] 1.2: Create `ActivityHandler` with `GET /api/activity` endpoint
  - Query params: `since` (timestamp), `limit` (default 50), `types` (comma-separated)
  - Fetch commits, agent states, pull statuses in parallel via goroutines
  - Merge and sort by timestamp descending
- [ ] 1.3: Add agent state change events
  - When agent starts/completes/fails, emit an event
  - Hook into existing `runAgentAsync` and `triggerAgent`
- [ ] 1.4: Wire pull status events
  - On `POST /api/pull` completion, emit pull event with repo/status

### Verification
- [x] Unit tests for event merging/sorting logic
- [x] Integration test for `/api/activity` endpoint
- [x] Build succeeds

---

## Phase 2: Frontend Timeline UI

### Tasks

- [x] 2.1: Create `static/activity.html` page
  - Mobile-first layout with timeline CSS
  - Event cards with type icon, repo, message, relative time
- [x] 2.2: Add `activity.js` client
  - Fetch `/api/activity` on load with pagination
  - Implement infinite scroll / load more
  - WebSocket subscription for real-time updates
- [x] 2.3: Add filter controls
  - Toggle buttons for commit/agent/pull event types
  - Persist filter state to localStorage
- [x] 2.4: Wire navigation
  - Add "Activity" link to header nav (app.js)

### Verification
- [ ] Page renders at localhost:activity
- [ ] Events display with correct formatting
- [ ] Filters toggle visibility correctly
- [ ] No console errors

---

## Phase 3: Real-Time Streaming

### Tasks

- [ ] 3.1: Extend WebSocket Hub to broadcast activity events
  - New `ActivityFeed` hub (or reuse existing `logHub`)
  - Broadcast event to all connected `/ws/activity` clients
- [ ] 3.2: Client subscribes to `/ws/activity`
  - On new event, prepend to timeline without reload
  - Deduplicate against existing items by event ID
- [ ] 3.3: Persist last-read cursor
  - Store `lastSeenEventID` in localStorage
  - On reconnect, fetch events since cursor

### Verification
- [ ] New events appear within 1s of occurring
- [ ] No duplicate events on reconnect
- [ ] Visual indicator for new events prepended

---

## Phase 4: Polish & Edge Cases

### Tasks

- [ ] 4.1: Empty state UI
  - "No activity yet" message with icon
- [ ] 4.2: Error handling
  - Show toast on fetch failure
  - Retry with exponential backoff
- [ ] 4.3: Performance optimization
  - Virtualized list for 500+ events
  - Debounce WebSocket renders
- [ ] 4.4: Final visual check and screenshot

### Verification
- [ ] Empty state renders correctly
- [ ] Errors show user-friendly messages
- [ ] Smooth scrolling with many events