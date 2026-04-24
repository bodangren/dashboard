# Lessons Learned

> This file is curated working memory, not an append-only log. Keep it at or below **50 lines**.
> Remove or condense entries that are no longer relevant to near-term planning.

## Architecture & Design

- (2026-03-28, git-view-enhance_20260328) CSS grid `auto-fill` with small minmax was too cramped at 3 cols for dense commit info. Fixed 2-col grid is better for terminal-styled dashboards with monospace text.
- (2026-03-29, agent-editor-fix_20260329) Section header comments in crontab are separate `Line` entries in the slice, not inline with the agent. When adding new agents with section headers, insert a `LineComment` before the `LineAgent`. When deleting, remove the preceding comment too.
- (2026-04-09, critical-bugs-rewrite_20260406) HandlerConfig pattern eliminated global mutable state in API handlers - pass dependencies via struct, not package-level vars. AgentID (schedule:directory:model) provides stable identity for crontab agents independent of array position.
- (2026-04-10, critical-bugs-rewrite_20260406) Adding `ToAPICommit()` method on git.Commit keeps packages decoupled while eliminating manual mapping loops in main.go.

## Recurring Gotchas

- (2026-03-28, git-view-enhance_20260328) Pure CSS/JS tracks don't have unit test coverage in a Go project. Manual verification is the only gate — always defer manual-smoke-test tasks until user can visually confirm.
- (2026-03-29, agent-editor-fix_20260329) Real crontab uses `>` (single) for redirect, not `>>`. Regex must handle both. OpenCode uses `-m` flag and `run <path>` positional, not `--model`/`--prompt`.
- (2026-04-09, critical-bugs-rewrite_20260406) Agent IDs with colons (:) must be URL-encoded when used in API paths. Use `url.PathEscape` not `url.QueryEscape` (latter encodes spaces as + which HTTP server doesn't decode back).
- (2026-04-10, critical-bugs-rewrite_20260406) When switching API endpoints, verify response structure differences — /api/repos returns {repos: [...]} wrapper while /api/projects returns an array directly.

## Patterns That Worked Well

- (2026-03-28, git-view-enhance_20260328) Keeping `.commit-age-badge` as a separate DOM element (not innerHTML string concat) made it easy to conditionally append only when commits exist.
- (2026-03-29, agent-editor-fix_20260329) Pending-comment pattern in parser — track the last comment seen, attach it to the next agent line, reset on non-comment/non-agent lines. Clean way to capture section headers without modifying the Line struct heavily.
- (2026-04-09, critical-bugs-rewrite_20260406) Harness detection via explicit name map {re *regexp.Regexp, name Harness} instead of deriving name from regex string slicing. Explicit is better than fragile string manipulation.
- (2026-04-10, critical-bugs-rewrite_20260406) nil vs empty slice in Go: `ReorganizeAutomation(nil)` treats all dirs as orphans, `ReorganizeAutomation([]string{})` processes normally. Use empty slice for consistent behavior.
- (2026-04-12, critical-bugs-rewrite_20260406) Mobile-first CSS: base styles for small screens, use `@media (min-width: 769px)` for desktop enhancements. Avoid chaining max-width queries; each builds on mobile, not overrides. Consolidate duplicate queries for same breakpoint into single @media block.
- (2026-04-13, critical-bugs-rewrite_20260406) Service worker API caching: sw.js should return early for `/api/` routes without hitting cache. Keep static asset caching separate from API network-only strategy.
- (2026-04-13, critical-bugs-rewrite_20260406) Shared utilities (esc, relativeTime) in utils.js work well across pages. Load utils.js before page-specific scripts to ensure functions available. Extract common functions early rather than duplicating across files.
- (2026-04-12, critical-bugs-rewrite_20260406) Agent timing visualization: show human-readable schedule (scheduleHuman) as primary, visual blocks as secondary detail. Add labels (day abbreviations, "Hours") for discoverability.
- (2026-04-12, critical-bugs-rewrite_20260406) CSS custom properties (--gap, --card-padding, --font-size-base) at :root enable easy responsive adjustments without hunting through multiple rules.
- (2026-04-17, coverage-improvement_20260417) ReadLogFile is testable with temp files; ReadCrontab/WriteCrontab exec crontab directly and require interface injection to test. Design functions to be testable from the start.

## Planning Improvements

- (2026-03-28, git-view-enhance_20260328) Responsive column count is subjective — spec said 3→2→1 but user found 3 too cramped. Should prototype layout decisions before committing to column counts in the spec.
- (2026-04-06, bugfix-three-bugs_20260406) Empty slice vs nil: Go's `[]string{}` is not nil. Tests asserting `nil` will fail if the function returns an empty initialized slice. Prefer returning `nil` for error/not-found paths.
- (2026-04-09, critical-bugs-rewrite_20260406) BUG-05 (stable agent IDs) required changes across backend (handler methods, AgentJSON struct) AND frontend (agents.js data-id, API calls). Plan such cross-cutting changes as a single atomic commit for easier rollback.
- (2026-04-10, critical-bugs-rewrite_20260406) Empty repos skipped via `err != nil || len(commits) == 0` — only skip on actual git errors, not zero commits.
- (2026-04-11, critical-bugs-rewrite_20260406) Frontend diff.js already expected author/message/timestamp fields before backend implemented them — suggests speculative frontend development or prior design discussion. Verify frontend expectations against actual API before implementing backend.
- (2026-04-12, critical-bugs-rewrite_20260406) Hardcoded binary path in JS form vs Go defaults: Go's buildAgentLine() already defaults to harness name when BinaryPath is empty. Frontend should pass empty string for new agents, Go handles the rest correctly.
- (2026-04-17, critical-bugs-rewrite_20260406) Many phases end with "manual verification" steps that can't be automated. Plan these as user-facing acceptance criteria rather than autonomous agent tasks. Code-complete is a valid track milestone.
- (2026-04-23, agent-orchestration-monitoring_20260423) WebSocket hub tests require real connections via httptest.NewServer — registering mockConn doesn't trigger ServeHTTP which populates clients. Use httptest.Server + websocket.DefaultDialer for integration tests.
- (2026-04-23, agent-orchestration-monitoring_20260423) conn.WriteJSON can panic on nil/closed websocket — mitigated by adding write deadline (1s) and removing Close() calls in error paths. Still needs robust error recovery for production use.
- (2026-04-23, agent-orchestration-monitoring_20260423) Write deadline on WebSocket WriteJSON prevents indefinite blocking; gorilla/websocket Close() can deadlock with server-side close — avoid calling Close() when connection state is uncertain.
- (2026-04-23, agent-orchestration-monitoring_20260423) Buffered channels (capacity 10) for register/unregister/broadcast prevent deadlocks in tests. Unbuffered channels cause blocking sends that can deadlock if goroutines aren't perfectly synchronized. Buffering also enables non-blocking broadcast behavior.
- (2026-04-23, agent-orchestration-monitoring_20260423) LogWatcher uses polling (500ms tick) to detect log file changes — simpler than inotify/fsnotify and works across platforms. Track file size to handle truncation.
- (2026-04-23, agent-orchestration-monitoring_20260423) WatcherManager StartWatching is idempotent — checks existence before creating new watcher. Prevents duplicate watchers for same agent.
- (2026-04-23, agent-orchestration-monitoring_20260423) Async agent execution via goroutine in triggerAgent — don't block HTTP response. nil Stdout/Stderr inherits parent fds; explicit pipes needed for WebSocket streaming.
- (2026-04-24, agent-orchestration-monitoring_20260423) Error state tracking: AgentStateMap stores per-agent state; stateMap passed into runAgentAsync; clear state on successful exit (0), set error state on non-zero exit. Frontend conditionally renders error badge when exit_code is present in API response.

## New Insights

- (2026-04-24, hub-panic-recovery_20260424) When testing panic recovery in goroutines that outlive the test, log.SetOutput(nil) causes panic in log output. Use log.SetOutput(oldWriter) with defer to restore, not nil. Save old writer before redirect to restore properly.
- (2026-04-24, hub-panic-recovery_20260424) WebSocket connections in tests can panic when hub tries to WriteJSON after test closes the connection. This is expected behavior — panic is recovered by Hub.run() and logged, not a test failure.
- (2026-04-24, hub-panic-recovery_20260424) log.SetOutput(&logBuf) redirects ALL log output including from other goroutines running concurrently. If you set log output to a buffer and another goroutine panics and logs, you may capture that output too.
- (2026-04-24, api-pull-status-endpoint_20260424) Adding tracking state (inProgress, lastPullTime, lastPullErr) to Handler struct requires initializing maps in NewHandler to avoid nil map panics on write.
- (2026-04-24, api-pull-status-endpoint_20260424) mutex unlock must happen in defer or after lock is held — always unlock before subsequent operations to avoid race conditions on map access.