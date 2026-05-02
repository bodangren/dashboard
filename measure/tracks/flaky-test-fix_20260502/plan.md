# Plan: Fix Flaky TestHub_Run_PanicRecoveryContainsMessage

## Context
The test at `internal/ws/hub_test.go:178` is flaky due to race between LogWatcher.Stop and conn.Read. The tech debt item FLaky-01 identifies this.

## Tasks

### Phase 1: Analyze the race condition
- [ ] Read hub.go and hub_test.go to understand LogWatcher.Stop and conn.Read interaction
- [ ] Run the test multiple times with -race to observe the failure pattern
- [ ] Identify where the data race occurs

### Phase 2: Implement fix
- [ ] Ensure LogWatcher.Stop() properly coordinates with connection reads
- [ ] Ensure conn.Read goroutine exits cleanly when stopped
- [ ] Verify with multiple runs of `go test -race`

### Phase 3: Verify
- [ ] Run test suite (go test -race ./internal/ws/...)
- [ ] Run all tests (go test ./...)
- [ ] Commit