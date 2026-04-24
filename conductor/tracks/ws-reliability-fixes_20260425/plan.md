# Implementation Plan: WebSocket Reliability & Security Fixes

## Phase 1: Fix AgentStateMap Race Condition (RACE-01)

### Tasks

- [ ] 1.1: Write failing test for concurrent AgentStateMap access
- [ ] 1.2: Add sync.RWMutex to AgentStateMap struct
- [ ] 1.3: Protect Get with RLock
- [ ] 1.4: Protect Set with Lock
- [ ] 1.5: Protect Clear with Lock
- [ ] 1.6: Run tests and verify fix

## Phase 2: Fix LogStreamHandler Goroutine Leak (LEAK-01)

### Tasks

- [ ] 2.1: Write test demonstrating goroutine leak
- [ ] 2.2: Modify ServeHTTP to use conn Read goroutine for disconnect detection
- [ ] 2.3: Use context cancellation or channel for exit signaling
- [ ] 2.4: Run tests and verify fix

## Phase 3: Fix Hub Broadcast FD Leak (LEAK-02)

### Tasks

- [ ] 3.1: Write test demonstrating FD leak
- [ ] 3.2: Add conn.Close() when removing from clients map
- [ ] 3.3: Add conn.Close() when removing from subscriptions
- [ ] 3.4: Run tests and verify fix

## Phase 4: Add Binary Path Validation (SAFE-01)

### Tasks

- [ ] 4.1: Write test for invalid binary path rejection
- [ ] 4.2: Create whitelist of allowed binaries
- [ ] 4.3: Add validation in triggerAgent before exec.Command
- [ ] 4.4: Run tests and verify fix

## Phase 5: Fix Nil ProcessState (SAFE-02)

### Tasks

- [ ] 5.1: Write test for nil ProcessState scenario
- [ ] 5.2: Add check for cmd.ProcessState != nil before ExitCode()
- [ ] 5.3: Run tests and verify fix

## Phase 6: Fix Hub.Stop() Channel Closure (LEAK-03)

### Tasks

- [ ] 6.1: Write test for Hub.Stop() behavior
- [ ] 6.2: Implement done channel pattern in Hub.run()
- [ ] 6.3: Handle channel closure without spinning
- [ ] 6.4: Run tests and verify fix

## Phase 7: Final Verification

### Tasks

- [ ] 7.1: Run full test suite (`go test ./...`)
- [ ] 7.2: Verify `go build` completes without errors
- [ ] 7.3: Update tech-debt.md with resolved items
- [ ] 7.4: Update lessons-learned.md with new insights
- [ ] 7.5: Commit with git note