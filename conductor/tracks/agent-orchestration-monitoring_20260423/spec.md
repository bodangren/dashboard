# Specification: Enhanced Agent Orchestration & Monitoring

## Overview

Extend the agent monitoring system with real-time log streaming via WebSockets, manual trigger overrides for cron jobs, and improved error reporting for failed agent runs.

## Functional Requirements

### 1. WebSocket Log Streaming
- Add `/ws/logs` WebSocket endpoint that streams agent log updates in real-time
- Frontend connects to WebSocket and receives live log entries as they are written
- Log stream includes: agent ID, timestamp, log line content, and entry type (stdout/stderr/info/error)
- Connection survives page navigation (use context to maintain WebSocket)
- Graceful degradation: if WebSocket unavailable, fall back to polling `/api/logs/<agentId>`

### 2. Manual Cron Trigger Override
- Add "Run Now" button to each agent card in the UI
- Clicking triggers `POST /api/agents/<agentId>/trigger` which executes the agent immediately (outside cron schedule)
- Triggered runs show visual indicator (spinner, "Running..." state) in the UI
- Multiple manual triggers are allowed (no locking)
- Manual triggers write to same log file as scheduled runs

### 3. Improved Error Reporting
- When agent exits with non-zero code, capture exit code and stderr output
- Store last error per agent in memory (not persisted, resets on server restart)
- Display error status on agent card: "Error (exit 1)" badge when last run failed
- Show last error message on hover or click of error badge
- Error badge clears when agent runs successfully

## Non-Functional Requirements

- WebSocket connections are read-only (server pushes, no client commands)
- WebSocket reconnection with exponential backoff on disconnect
- Manual trigger should not block HTTP response (run async)
- All existing tests continue to pass

## Acceptance Criteria

1. WebSocket endpoint `/ws/logs` accepts connections and streams log data
2. Frontend displays real-time log lines in a scrollable panel when viewing agent details
3. "Run Now" button appears on each agent card
4. Clicking "Run Now" triggers immediate execution, shows running state, updates log
5. Failed agents show error badge with exit code
6. Error badge shows stderr content on hover/click
7. `go test ./...` passes with >80% coverage on agents and api packages
8. `go build` completes without errors

## Out of Scope

- Persisting error history across server restarts
- Agent cancellation/termination
- Authentication (single-user local app)
- Push notifications