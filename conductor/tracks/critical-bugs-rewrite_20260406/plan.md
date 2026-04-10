# Implementation Plan: Critical Bugs & Display Rewrite

## Phase 1 — Backend Bug Fixes (P0)

### 1.1 BUG-01: Eliminate global mutable state in API handlers

- [ ] 1.1.1 Write test: `TestNewHandlerConfig_NoGlobals` — create handler via config struct, verify functions work without any package-level vars
- [ ] 1.1.2 Create `HandlerConfig` struct in `handlers.go` with `GetCommitsFunc`, `GetDiffFunc`, `PullFunc` fields
- [ ] 1.1.3 Update `RegisterRoutes` to accept `HandlerConfig` and store funcs in handler fields directly, remove `defaultGetCommits`, `defaultGetDiff`, `defaultPullRepo` package-level vars
- [ ] 1.1.4 Remove `SetGitFuncs`, `SetPullFunc` functions from `handlers.go`
- [ ] 1.1.5 Update `main.go` to construct `HandlerConfig` and pass to `RegisterRoutes`
- [ ] 1.1.6 Update `integration_test.go` to use the new constructor pattern instead of mutating globals
- [ ] 1.1.7 Update `handlers_test.go` helper `newTestHandler` / `newTestHandlerWithPull` to use `HandlerConfig`
- [ ] 1.1.8 Run `go test -race ./...` — confirm no data races
- [ ] 1.1.9 Conductor - User Manual Verification 'Phase 1.1' (Protocol in workflow.md)

### 1.2 BUG-02: Fix `detectHarness` fragile regex slicing

- [ ] 1.2.1 Write test: `TestDetectHarnessExplicitMap` — verify each harness name is returned correctly regardless of regex pattern syntax
- [ ] 1.2.2 Replace `detectHarness` with a loop over `[]struct{re *regexp.Regexp; name Harness}` so harness name is explicit, not derived from regex source string
- [ ] 1.2.3 Run `go test ./internal/agents/...` — all parser/writer tests pass
- [ ] 1.2.4 Conductor - User Manual Verification 'Phase 1.2' (Protocol in workflow.md)

### 1.3 BUG-03: Fix `isEnvVarLine` over-matching

- [ ] 1.3.1 Write test: `TestIsEnvVarLine_DoesNotMatchCronWithEquals` — input `0 */4 * * * cd /foo && make test=1` returns false
- [ ] 1.3.2 Write test: `TestIsEnvVarLine_MatchesRealEnvVar` — input `SHELL=/bin/bash` returns true, `PATH=/usr/bin:/bin` returns true
- [ ] 1.3.3 Simplify `isEnvVarLine` to use only the regex `^[A-Za-z_][A-Za-z0-9_]*=` (compile once at package level), remove the `strings.Contains` clause
- [ ] 1.3.4 Run `go test ./internal/agents/...` — all tests pass
- [ ] 1.3.5 Conductor - User Manual Verification 'Phase 1.3' (Protocol in workflow.md)

### 1.4 BUG-04: Fix `toggleAgent` fragile re-lookup

- [ ] 1.4.1 Write test: `TestToggleAgentHandler_ReturnsUpdatedState` — toggle agent, verify response has correct `enabled` value
- [ ] 1.4.2 Simplify `toggleAgent`: after toggle+write, re-read crontab, find agent by matching schedule+directory+model, return it. Remove the confusing `||` logic and nil fallback
- [ ] 1.4.3 Run `go test ./internal/api/...` — all tests pass
- [ ] 1.4.4 Conductor - User Manual Verification 'Phase 1.4' (Protocol in workflow.md)

### 1.5 BUG-05: Replace index-based agent referencing with stable IDs

- [ ] 1.5.1 Write test: `TestAgentID_IsDeterministic` — same schedule+directory+model produces same ID
- [ ] 1.5.2 Write test: `TestAgentID_DifferentForDifferentAgents` — different schedules produce different IDs
- [ ] 1.5.3 Add `AgentID()` function to `agents` package: returns `fmt.Sprintf("%s:%s:%s", a.Schedule, a.Directory, a.Model)`
- [ ] 1.5.4 Add `ID string` field to `AgentJSON` and populate it in `agentToJSON`
- [ ] 1.5.5 Update `AgentHandler.resolveCrontabAndAgent` to accept an `id string` and find agent by ID instead of numeric index
- [ ] 1.5.6 Update all handler methods (`deleteAgent`, `toggleAgent`, `updateAgent`, `getLog`) to resolve by ID
- [ ] 1.5.7 Update `static/agents.js`: send agent `id` instead of `idx` in all API calls
- [ ] 1.5.8 Write test: `TestAgentDeleteByID` — delete agent by ID, verify correct agent removed
- [ ] 1.5.9 Run `go test ./...` — all tests pass
- [ ] 1.5.10 Conductor - User Manual Verification 'Phase 1.5' (Protocol in workflow.md)

---

## Phase 2 — Architecture & Data Layer Cleanup

### 2.1 ARCH-01: Eliminate duplicate Commit type

- [x] 2.1.1 Add `ToAPICommit()` method to `git.Commit` that returns `api.Commit`
- [x] 2.1.2 Simplify `main.go` to use `git.Commit.ToAPICommit()` instead of manual field mapping loop
- [x] 2.1.3 Run `go test ./...` — all tests pass
- [x] 2.1.4 Conductor - User Manual Verification 'Phase 2.1' (Protocol in workflow.md)

### 2.2 ARCH-02: Add lightweight `/api/repos` endpoint

- [x] 2.2.1 Write test: `TestReposHandler_returnsOnlyNamesAndPaths` — verify response has `name` and `path` fields, no `commits` or `last_commit_at`
- [x] 2.2.2 Add `repos` handler method to `Handler` that returns `[{name, path}]` from `h.repos` without calling `getCommits`
- [x] 2.2.3 Register `/api/repos` route in `RegisterRoutes`
- [x] 2.2.4 Update `agents.js` `fetchRepos` to use `/api/repos` instead of `/api/projects`
- [x] 2.2.5 Run `go test ./internal/api/...` — all tests pass
- [x] 2.2.6 Conductor - User Manual Verification 'Phase 2.2' (Protocol in workflow.md)

### 2.3 ARCH-03: Show empty repos in the dashboard

- [x] 2.3.1 Write test: `TestProjectsHandler_includesEmptyRepos` — repo with zero commits appears in response with empty `commits` array
- [x] 2.3.2 Update `projects` handler: only skip repo on `err != nil` (actual git error), include repos with zero commits
- [x] 2.3.3 Update `renderProject` in `app.js`: show "no commits yet" message when `project.commits.length === 0`
- [x] 2.3.4 Run `go test ./internal/api/...` — all tests pass
- [x] 2.3.5 Conductor - User Manual Verification 'Phase 2.3' (Protocol in workflow.md)

### 2.4 ARCH-04: Remove API caching from Service Worker

- [x] 2.4.1 Update `sw.js`: remove the `if (url.pathname.startsWith('/api/'))` block entirely; API requests go straight to network with no caching
- [ ] 2.4.2 Manual test: load dashboard, go offline in DevTools, verify API requests fail cleanly (no stale cache served)
- [ ] 2.4.3 Conductor - User Manual Verification 'Phase 2.4' (Protocol in workflow.md)

### 2.5 ARCH-05: Fix `AddAgent` standalone `ReorganizeAutomation(nil)` call

- [x] 2.5.1 Write test: `TestAddAgentWithoutProjectList` — add agent without project list, verify it appears in output
- [x] 2.5.2 Update `AddAgent` to call `ReorganizeAutomation` with empty project list instead of nil (consistent behavior)
- [x] 2.5.3 Run `go test ./internal/agents/...` — all tests pass
- [x] 2.5.4 Conductor - User Manual Verification 'Phase 2.5' (Protocol in workflow.md)

---

## Phase 3 — Frontend Display Overhaul

### 3.1 UX-05: Extract shared utilities into `utils.js`

- [ ] 3.1.1 Create `static/utils.js` containing `esc()` and `relativeTime()` (with future-timestamp fix from UX-06)
- [ ] 3.1.2 Remove `esc()` from `app.js`, `agents.js`, `diff.js`
- [ ] 3.1.3 Remove `relativeTime` from `app.js`, add import comment noting it comes from `utils.js`
- [ ] 3.1.4 Add `<script src="utils.js"></script>` to `index.html`, `agents.html`, `diff.html` before page-specific scripts
- [ ] 3.1.5 Manual test: load all three pages, verify no console errors, all features work
- [ ] 3.1.6 Conductor - User Manual Verification 'Phase 3.1' (Protocol in workflow.md)

### 3.2 UX-06: Fix `relativeTime` for future timestamps

- [ ] 3.2.1 Update `relativeTime` in `utils.js`: if `diffMs < 0`, return "just now"
- [ ] 3.2.2 Manual test: verify no negative time displays on any commit
- [ ] 3.2.3 Conductor - User Manual Verification 'Phase 3.2' (Protocol in workflow.md)

### 3.3 UX-01: Clean up project card display

- [ ] 3.3.1 Update `renderProject` in `app.js`: show project name only, add full path as `title` attribute on the name element for hover tooltip
- [ ] 3.3.2 Update `style.css`: hide `.project-path` display by default (or remove it from render entirely, keep only as tooltip)
- [ ] 3.3.3 Manual test: verify project cards show clean names, paths visible on hover only
- [ ] 3.3.4 Conductor - User Manual Verification 'Phase 3.3' (Protocol in workflow.md)

### 3.4 UX-04: Add commit metadata to diff page

- [ ] 3.4.1 Write test: `TestDiffHandler_returnsMetadata` — verify response includes `message`, `author`, `timestamp` fields
- [ ] 3.4.2 Add `GetCommitInfoFunc` signature: `func(repoPath, hash string) (message, author string, timestamp time.Time, err error)`
- [ ] 3.4.3 Implement in `git` package: run `git log -1 --format="%s%x1f%an%x1f%ct" <hash>` and parse result
- [ ] 3.4.4 Extend `DiffResponse` to include `Message`, `Author`, `Timestamp` fields
- [ ] 3.4.5 Update `diff` handler to call `getCommitInfo` alongside `getDiff`, populate response
- [ ] 3.4.6 Update `static/diff.js`: set `titleEl.textContent` to `data.message || data.hash`, populate `metaEl` with author and formatted date
- [ ] 3.4.7 Run `go test ./internal/api/...` — new tests pass
- [ ] 3.4.8 Manual test: click a commit, verify diff page shows message as title, author and date in meta section
- [ ] 3.4.9 Conductor - User Manual Verification 'Phase 3.4' (Protocol in workflow.md)

### 3.5 UX-02: Remove hardcoded binary path from agent form

- [ ] 3.5.1 Update `buildForm` in `agents.js`: default binary path to empty string. Let `buildAgentLine` in Go default to harness name when `BinaryPath` is empty. Only use stored `binary_path` when editing an existing agent
- [ ] 3.5.2 Manual test: open agent create form, verify binary path is not hardcoded to nvm path
- [ ] 3.5.3 Conductor - User Manual Verification 'Phase 3.5' (Protocol in workflow.md)

### 3.6 UX-03: Improve agent timing visualization

- [ ] 3.6.1 Update `renderTimingVisualization` in `agents.js`: add day labels ("S M T W T F S") as small text above day blocks, add "Hours" label before hour blocks
- [ ] 3.6.2 Show `scheduleHuman(cron)` output as primary readable text, with visual blocks as secondary detail
- [ ] 3.6.3 Add small legend CSS in `style.css`: `.sched-legend` class with tiny text showing green=active, red=inactive
- [ ] 3.6.4 Manual test: verify agent cards have labeled timing visualization with readable schedule
- [ ] 3.6.5 Conductor - User Manual Verification 'Phase 3.6' (Protocol in workflow.md)

### 3.7 UX-07: Refactor CSS to mobile-first

- [ ] 3.7.1 Rewrite `style.css` base rules as mobile (small) styles
- [ ] 3.7.2 Replace all `@media (max-width: 768px)` overrides with `@media (min-width: 768px)` for desktop enhancements
- [ ] 3.7.3 Consolidate duplicate media queries for the same breakpoint into single `@media` blocks
- [ ] 3.7.4 Introduce CSS custom properties for responsive values: `--gap`, `--card-padding`, `--font-size-base`
- [ ] 3.7.5 Manual test: verify dashboard renders correctly at 375px (mobile), 768px (tablet), 1400px (desktop)
- [ ] 3.7.6 Conductor - User Manual Verification 'Phase 3.7' (Protocol in workflow.md)

---

## Phase 4 — Final Verification & Cleanup

### 4.1 Full regression

- [ ] 4.1.1 Run `go test -race ./...` — zero races, all tests pass
- [ ] 4.1.2 Run `go vet ./...` — no warnings
- [ ] 4.1.3 Build: `go build -o dashboard .` — succeeds
- [ ] 4.1.4 Manual smoke test: load dashboard, verify all repos appear (including empty ones), click commit to view diff with message/author/date, navigate to agents page and verify agents load fast without git log overhead, create/edit/toggle/delete an agent, verify timing visualization has labels
- [ ] 4.1.5 Conductor - User Manual Verification 'Phase 4' (Protocol in workflow.md)

### 4.2 Update project memory

- [ ] 4.2.1 Update `conductor/lessons-learned.md` with insights from this track
- [ ] 4.2.2 Update `conductor/tech-debt.md`: mark resolved items, add any new shortcuts discovered
- [ ] 4.2.3 Conductor - User Manual Verification 'Phase 4.2' (Protocol in workflow.md)

### 4.3 Final commit

- [ ] 4.3.1 Commit: `fix(all): resolve critical bugs, rewrite display layer, eliminate data races`
