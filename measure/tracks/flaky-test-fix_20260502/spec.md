# Track: Fix Flaky TestHub_Run_PanicRecoveryContainsMessage

## Problem
`TestHub_Run_PanicRecoveryContainsMessage` fails non-deterministically. The race detector shows a data race between `LogWatcher.Stop` and `conn.Read`.

## Root Cause
The test creates a Hub with LogWatcher, registers a client connection, triggers Hub.run(), and then the main goroutine calls `w.Stop()`. But `conn.Read()` is also running concurrently in another goroutine. The race is between Stop closing the connection and Read continuing.

## Solution
The LogWatcher.Stop() method needs to properly synchronize with any in-flight connection reads. The connection close should happen once, and Read goroutines should exit cleanly when the connection is closed.

## Acceptance Criteria
1. TestHub_Run_PanicRecoveryContainsMessage passes consistently with -race flag
2. LogWatcher.Stop() properly signals all goroutines to exit before closing connections
3. No data races detected by go test -race