# Specification: WebSocket Reliability & Security Fixes

## Overview

Fix critical security and reliability issues identified in the WebSocket hub, agent executor, and log streaming components. These issues were discovered during code review and documented in the tech-debt registry.

## Issues to Address

### 1. RACE-01: AgentStateMap Thread Safety (High)
**Location:** `internal/agents/types.go:104-117`
**Problem:** AgentStateMap has no mutex — concurrent goroutine writes from runAgentAsync vs HTTP handler reads.
**Fix:** Add sync.RWMutex to protect Get/Set/Clear operations.

### 2. LEAK-01: LogStreamHandler Goroutine Leak (High)
**Location:** `internal/ws/log_stream_handler.go:46`
**Problem:** `<-make(chan struct{})` never detects client disconnect, causing goroutine leak.
**Fix:** Use conn.Read goroutine to detect close, then signal exit.

### 3. LEAK-02: Hub Broadcast FD Leak (High)
**Location:** `internal/ws/hub.go:84-114`
**Problem:** Hub broadcast removes conns without calling conn.Close() — WebSocket FD leak.
**Fix:** Add conn.Close() when removing from clients/subscriptions.

### 4. SAFE-01: Binary Path Validation (Critical)
**Location:** `internal/api/agent_handlers.go:257-321`
**Problem:** triggerAgent runs exec.Command with user-supplied binary path — no validation.
**Fix:** Add whitelist validation for BinaryPath/Harness.

### 5. SAFE-02: Nil ProcessState Dereference (Critical)
**Location:** `internal/ws/agent_executor.go:49`
**Problem:** agent_executor.go nil dereference on cmd.ProcessState after unexpected kill.
**Fix:** Check ProcessState != nil before calling ExitCode().

### 6. LEAK-03: Hub.Stop() Channel Closure Spin (Medium)
**Location:** `internal/ws/hub.go:134-145`
**Problem:** Hub.Stop() closes channels causing run() to spin on zero values.
**Fix:** Use done channel pattern or check channel closure.

## Functional Requirements

1. All concurrent access to AgentStateMap must be thread-safe
2. LogStreamHandler must properly clean up goroutines on client disconnect
3. Hub must close WebSocket connections when removing them
4. Agent execution must validate binary paths against an allowed list
5. Agent executor must handle nil ProcessState gracefully
6. Hub.run() must not spin after channels are closed

## Acceptance Criteria

- [ ] AgentStateMap operations are protected by RWMutex
- [ ] LogStreamHandler goroutine exits when client disconnects
- [ ] Hub properly closes connections when removing them
- [ ] triggerAgent validates binary path before execution
- [ ] Agent executor handles nil ProcessState without panic
- [ ] Hub.run() handles channel closure gracefully
- [ ] All existing tests pass
- [ ] New tests cover the fixes