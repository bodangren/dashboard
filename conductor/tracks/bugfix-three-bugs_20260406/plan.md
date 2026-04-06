# Implementation Plan: Three Critical Bug Fixes

## Phase 1 — Fix Model Discovery (BUG-1) ✓

- [x] 1.1.1 Write test: `TestDiscoverModels_WithExplicitPath`
- [x] 1.1.2 Update `discoverModels()` to accept optional binary path with `resolveOpenCodeBinary()` fallback
- [x] 1.1.3 Update `HandleModels` to use `ah.openCodeBin`
- [x] 1.1.4 Update `main.go` to pass empty string (auto-resolve) to `AgentHandler`
- [x] 1.1.5 Run `go test ./internal/api/...` — all tests pass
- [x] 1.1.6 Manual test: `curl /api/models` returns 252 models entries

- [ ] 1.1.7 Conductor - User Manual Verification 'Phase 1.1' (Protocol in workflow.md)

## Phase 2 — Fix Pull Error Reporting (BUG-2) ✓
- [x] 2.1.1 Update pull handler to return JSON `{"status": "error", "error": "..."}``
- [x] 2.1.2 Write test: `TestPullHandler_pullError` — verify JSON error response
- [x] 2.1.3 Update `app.js` pull button to show actual git error message from JSON
- [x] 2.1.4 Run `go test ./internal/api/...` — all tests pass
- [ ] 2.1.5 Manual test: pull a repo with no tracking branch, verify UI shows meaningful error
- [ ] 2.1.6 Conductor - User Manual Verification 'Phase 2.1' (Protocol in workflow.md)
- [ ] 2.2 Add pull status summary endpoint (deferred to future track)

- [ ] 2.2.1 Add `/api/pull/status` GET endpoint
- [ ] 2.2.2 Conductor - User Manual Verification 'Phase 2.2' (Protocol in workflow.md)

## Phase 3 — Fix Agent Log Display (BUG-3) ✓
- [x] 3.1.1 Write test: `TestGetLog_ResolvesRelativePath`
- [x] 3.1.2 Write test: `TestGetLog_AbsolutePathUnchanged`
- [x] 3.1.3 Update `getLog` to resolve relative paths and expand `$HOME`/`~`
- [x] 3.1.4 Run `go test ./internal/api/...` — all tests pass
- [x] 3.1.5 Manual test: verified logs display for agent cards with relative paths
- [x] 3.1.6 Conductor - User Manual Verification 'Phase 3.1' (Protocol in workflow.md)

- [x] 3.2.1 Write test: `TestGetLog_ExpandsHomeEnvVar`
- [x] 3.2.2 Run `go test ./internal/api/...` — all tests pass
- [x] 3.2.3 Conductor - User Manual Verification 'Phase 3.2' (Protocol in workflow.md)

## Phase 4 — Final Verification & Cleanup
- [x] 4.1.1 Run `go test -race ./...` — zero races, all tests pass (a17e92e)
- [x] 4.1.2 Run `go vet ./...` — no warnings (a17e92e)
- [x] 4.1.3 Build: `go build -o dashboard .` — succeeds (a17e92e)
- [ ] 4.1.4 Manual smoke test: verify models, pulls, logs
 agents page
- [ ] 4.1.5 Conductor - User Manual Verification 'Phase 4' (Protocol in workflow.md)
- [x] 4.2.1 Update `conductor/lessons-learned.md` (a17e92e)
- [x] 4.2.2 Update `conductor/tech-debt.md` (a17e92e)
- [ ] 4.2.3 Conductor - User Manual Verification 'Phase 4.2' (Protocol in workflow.md)
- [~] 4.3.1 Commit: `fix(models,pull,logs): resolve three critical bugs — model discovery, pull errors, relative log paths`
