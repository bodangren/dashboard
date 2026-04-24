# Implementation Plan: Agent Log Streaming via WebSocket

## Phase 1: Investigation & Design

### Tasks

- [x] 1.1: Examine current LogWatcher implementation in internal/ws/log_watcher.go
- [x] 1.2: Examine existing WebSocket hub structure in internal/ws/hub.go
- [x] 1.3: Design WebSocket endpoint route and message format

## Phase 2: Implementation

### Tasks

- [x] 2.1: Extend Hub to support per-agent subscriptions (targeted broadcast)
- [x] 2.2: Replace LogWatcher 500ms polling with inotify-based file watch
- [x] 2.3: Implement /ws/logs handler that accepts agent ID and streams logs
- [x] 2.4: Handle client disconnect and per-agent subscription cleanup

## Phase 3: Tests

### Tasks

- [x] 3.1: Write unit test for targeted broadcast and per-agent filtering
- [x] 3.2: Write unit test for inotify-based LogWatcher
- [x] 3.3: Run full test suite (`go test ./...`) — all tests must pass
- [x] 3.4: Verify `go build` completes without errors

## Phase 4: Finalize

### Tasks

- [x] 4.1: Update tech-debt.md with any new insights
- [x] 4.2: Update lessons-learned.md with any new insights
- [x] 4.3: Commit with git note