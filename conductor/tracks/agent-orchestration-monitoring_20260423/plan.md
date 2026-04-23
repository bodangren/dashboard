# Implementation Plan: Enhanced Agent Orchestration & Monitoring

## Phase 1: WebSocket Log Streaming Infrastructure

### Tasks

- [x] 1.1: Write unit tests for WebSocket hub (goroutine-safe broadcast, client registration/cleanup) [cf82641]
- [x] 1.2: Implement WebSocket hub struct with client tracking and broadcast channel [cf82641]
- [x] 1.3: Write unit tests for `/ws/logs` HTTP upgrade handler
- [x] 1.4: Implement `/ws/logs` WebSocket endpoint with goroutine-safe client management
- [x] 1.5: Write unit tests for log broadcaster (buffered channel, non-blocking send)
- [x] 1.5: Write unit tests for log broadcaster (buffered channel, non-blocking send)
- [x] 1.6: Integrate log broadcaster into agent execution to stream stdout/stderr

[checkpoint: a128d56]

## Phase 2: Manual Cron Trigger Override

### Tasks

- [x] 2.1: Write unit tests for agent trigger handler (POST /api/agents/<id>/trigger)
- [x] 2.2: Implement trigger endpoint that queues agent for immediate execution
- [x] 2.3: Write integration test for concurrent manual trigger handling
- [x] 2.4: Add "Run Now" button to agent card in agents.js
- [x] 2.5: Implement running state UI (spinner, "Running..." badge) on agent card

## Phase 3: Improved Error Reporting

### Tasks

- [x] 3.1: Write unit tests for agent error state tracking (lastError field, exitCode capture)
- [x] 3.2: Implement error capture on non-zero agent exit
- [x] 3.3: Write unit tests for error display in API response
- [x] 3.4: Add error badge to agent card in agents.js (show on failed state)
- [x] 3.5: Implement hover/click tooltip showing last error stderr

## Phase 4: Integration & Verification

### Tasks

- [x] 4.1: Run full test suite (`go test ./...`) — all tests must pass
- [x] 4.2: Run coverage check — agents and api packages >80%
- [x] 4.3: Verify `go build` completes without errors
- [P] 4.4: Manual verification of WebSocket streaming in browser
- [P] 4.5: Manual verification of "Run Now" trigger and running state
- [P] 4.6: Manual verification of error badge display