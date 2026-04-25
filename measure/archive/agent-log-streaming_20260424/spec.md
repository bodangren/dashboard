# Specification: Agent Log Streaming via WebSocket

## Overview

The dashboard currently polls agent log files via the LogWatcher which pushes updates to WebSocket clients. This track replaces the polling-based log watching with true real-time streaming that pushes new log lines to connected clients as they are written, without periodic polling intervals.

## Functional Requirements

1. **WebSocket Endpoint**: Expose a `/ws/logs` endpoint that accepts an agent ID and streams log updates in real-time
2. **Log Streaming**: When a client connects to `/ws/logs?agent=<encoded-agent-id>`, begin streaming log content as it is written
3. **Incremental Updates**: Send only new log lines since connection start, not entire file contents
4. **Fallback**: If no log file exists for the agent, send an appropriate error message and close
5. **Connection Lifecycle**: Handle client disconnect gracefully, clean up any subscription state

## Non-Functional Requirements

- **Low Latency**: New log lines should appear on the client within 100ms of being written
- **No Polling**: Use file system notifications or equivalent to detect changes, not periodic polling
- **Multiple Clients**: Multiple dashboard clients can connect to the same agent log stream simultaneously
- **Backpressure**: If a client is slow to consume messages, buffer appropriately but don't block logging

## Acceptance Criteria

1. GET `/ws/logs?agent=<encoded-id>` establishes a WebSocket connection
2. Log lines written after connection are pushed to the client within 100ms
3. Client receives `{type: "log", content: "...", timestamp: "..."}` messages
4. If agent doesn't exist or has no log file, connection is closed with error
5. Disconnecting client doesn't affect other connections or logging behavior
6. `go test ./...` passes
7. `go build` completes without errors

## Out of Scope

- Authentication for WebSocket connections (single-user local app)
- TLS for WebSocket (localhost only)
- Historical log retrieval (only streaming live logs)
- Log file rotation handling