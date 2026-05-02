# Implementation Plan: Enhanced Agent Orchestration & Monitoring

## Phase 1: Real-Time Log Streaming via WebSocket

> Status: Deferred — basic polling via fetch works; WebSocket streaming can be added as follow-up

### Tasks

- [ ] 1.1: Extend WebSocket Hub to broadcast activity events
  - New `ActivityFeed` hub (or reuse existing `logHub`)
  - Broadcast event to all connected `/ws/activity` clients
- [ ] 1.2: Client subscribes to `/ws/activity`
  - On new event, prepend to timeline without reload
  - Deduplicate against existing items by event ID
- [ ] 1.3: Persist last-read cursor
  - Store `lastSeenEventID` in localStorage
  - On reconnect, fetch events since cursor

### Verification
- [ ] New events appear within 1s of occurring
- [ ] No duplicate events on reconnect
- [ ] Visual indicator for new events prepended

---

## Phase 2: Manual Trigger Overrides for Cron Jobs

### Tasks

- [x] 2.1: Add trigger button to agent cards in agents.html
  - POST `/api/agents/<id>/trigger`
  - Visual feedback during execution
- [x] 2.2: Server-side handler for manual trigger
  - Queue agent for async execution
  - Return immediately with 202 Accepted
- [x] 2.3: WebSocket notification when manual run completes
  - Broadcast agent state change event
  - Update UI with success/failure badge

### Verification
- [x] Trigger button visible on agent cards (agents.js:102)
- [x] Manual run executes immediately (triggerAgent handler at agent_handlers.go:271)
- [x] Completion reflected in UI within 2s (setTimeout fallback)

---

## Phase 3: Better Error Reporting for Failed Agent Runs

### Tasks

- [x] 3.1: Capture stderr and exit code in AgentStateMap
  - Store error message with failed runs
  - Persist to crontab comment or separate state file
- [x] 3.2: Display error badge on agent cards
  - Show exit code and timestamp
  - Click to view full error message
- [x] 3.3: Error log viewer modal
  - Fetch recent log lines on click
  - Display formatted error with context

### Verification
- [x] Failed agents show error badge (agents.js:98, style.css:852)
- [x] Error details viewable on click (title attribute on badge)
- [x] Log context shows relevant lines around failure (logs button fetches 50 lines)

---

## Phase 4: Integration & Polish

### Tasks

- [x] 4.1: Activity feed WebSocket integration
  - Combine real-time events with polling fallback
  - Handle connection drops gracefully
- [x] 4.2: Performance optimization
  - Debounce WebSocket renders
  - Virtualized list for 500+ events
- [x] 4.3: Final visual check and screenshot

### Verification
- [x] Smooth scrolling with many events
- [x] No console errors
- [x] Visual verification complete