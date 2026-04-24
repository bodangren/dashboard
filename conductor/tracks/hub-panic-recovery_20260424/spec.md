# Specification: Hub.run() Panic Recovery

## Overview

The WebSocket hub's `run()` method currently recovers from panics with an empty `recover()` call that silently discards panic information. This makes debugging production issues difficult since panics are swallowed without any logging or visibility.

## Functional Requirements

1. **Panic Logging**: When a goroutine panics, recover the panic value and log it with appropriate context (agent ID, connection info if available)
2. **Debuggability**: Log should include the panic message, stack trace if possible, and enough context to identify which client/agent was affected
3. **Non-blocking**: Panic recovery must not block the hub's event loop — other clients must continue to be served

## Non-Functional Requirements

- Log format should be consistent with existing logging patterns
- Must not introduce new synchronization issues (panic + logging could deadlock if locks are held)
- Keep recovery in `run()` central — do not scatter recovery logic throughout

## Acceptance Criteria

1. `Hub.run()` recovers panics and logs the panic value and stack
2. Panic recovery does not block other clients
3. Log output includes enough context to identify the affected connection
4. `go test ./...` continues to pass
5. `go build` completes without errors