# Plan: Fix Flaky TestHub_Run_PanicRecoveryContainsMessage

## Context
The test at `internal/ws/hub_test.go:178` is flaky due to race between LogWatcher.Stop and conn.Read. The tech debt item FLaky-01 identifies this.

## Tasks

### Phase 1: Analyze the race condition
- [x] Read hub.go and hub_test.go to understand LogWatcher.Stop and conn.Read interaction
- [x] Run the test multiple times with -race to observe the failure pattern
- [x] Identify where the data race occurs

### Phase 2: Implement fix
- [x] Ensure LogWatcher.Stop() properly coordinates with connection reads
- [x] Ensure conn.Read goroutine exits cleanly when stopped
- [x] Verify with multiple runs of `go test -race`

### Phase 3: Verify
- [x] Run test suite (go test -race ./internal/ws/...)
- [x] Run all tests (go test ./...)
- [x] Commit