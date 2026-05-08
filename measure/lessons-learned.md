# Lessons Learned

> This file is curated working memory, not an append-only log. Keep it at or below **50 lines**.
> Remove or condense entries that are no longer relevant to near-term planning.

## Architecture & Design

- (2026-03-29, agent-editor-fix_20260329) Section header comments in crontab are separate `Line` entries in the slice, not inline with the agent. When adding new agents with section headers, insert a `LineComment` before the `LineAgent`. When deleting, remove the preceding comment too.
- (2026-04-09, critical-bugs-rewrite_20260406) HandlerConfig pattern eliminated global mutable state in API handlers - pass dependencies via struct, not package-level vars. AgentID (schedule:directory:model) provides stable identity for crontab agents independent of array position.
- (2026-04-10, critical-bugs-rewrite_20260406) Adding `ToAPICommit()` method on git.Commit keeps packages decoupled.

## Recurring Gotchas

- (2026-03-28/29, git-view-enhance_20260328, agent-editor-fix_20260329) Pure CSS/JS tracks lack unit test coverage in Go projects; manual verification only. Real crontab uses `>` (single) for redirect, not `>>`; OpenCode uses `-m` flag and `run <path>` positional, not `--model`/`--prompt`.
- (2026-04-09, critical-bugs-rewrite_20260406) Agent IDs with colons (:) must be URL-encoded when used in API paths. Use `url.PathEscape` not `url.QueryEscape` (latter encodes spaces as + which HTTP server doesn't decode back).

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

- (2026-04-09, critical-bugs-rewrite_20260406) Plan cross-cutting changes as a single atomic commit for easier rollback.
- (2026-04-17, critical-bugs-rewrite_20260406) Manual verification steps can't be automated. Plan as user-facing acceptance criteria.
- (2026-04-23, agent-orchestration-monitoring_20260423) WebSocket hub tests require real connections via httptest.NewServer + websocket.DefaultDialer.
- (2026-04-23, agent-orchestration-monitoring_20260423) Buffered channels (capacity 10) for register/unregister/broadcast prevent deadlocks.

## New Insights

- (2026-04-24, hub-panic-recovery_20260424) When testing panic recovery in goroutines that outlive the test, log.SetOutput(nil) causes panic. Use log.SetOutput(oldWriter) with defer to restore properly. log.SetOutput(&logBuf) captures ALL concurrent log output including from other panicking goroutines.
- (2026-05-02, commit-ai-insights_20260502) mockProvider in ai package uses same interface as real LLM provider, allowing easy test substitution. Cache summaries with TTL to avoid redundant LLM calls. Detect flags (conflict-markers, WIP, rapid-changes) from commit message/body to highlight issues. ActivityEnhancer ties summarizer to activity handler, enriching events on demand.
- (2026-05-02, enhanced-agent-orchestration_20260502) ActivityHub broadcasts to all connected clients; when agents trigger, RecordAgentEvent fires both the recentAgentEvents slice (for /api/activity) and broadcasts via activityHub to WebSocket subscribers for real-time updates.
- (2026-05-03, commit-ai-insights_20260502) To avoid import cycles, define local types (CommitInfo) that mirror external types (git.Commit) rather than importing the external package. ActivityEnhancer takes []CommitInfo instead of []git.Commit to keep ai package independent.
- (2026-05-03, repo_health_monitoring_20260503) Health data from /api/projects (health field in Project struct); status: clean=healthy, minor=warning, significant=critical. GetStaleBranches uses git for-each-ref with committerdate; empty repos have no branches.
- (2026-05-05, platform-git-features_20260505) Git format strings with | as delimiter work for parsing multiple fields. Stash@{N} refs must strip `stash@{` prefix and `}` suffix. Branch ops in tests require repo with at least one commit.
- (2026-05-04, dashboard_customization_20260503) Theme system using CSS custom properties on :root works well; inline script in <head> prevents flash. CSS transitions (0.2s ease) on background-color/border-color feel polished. Ctrl+T for theme toggle enhances power users.
- (2026-05-06, flaky-test-fix_20260502) sync.WaitGroup better than channels for test goroutine coordination — avoids logBuf race. sync.Mutex protects inotifyFd writes during initInotify vs Stop.
- (2026-05-06, notification_alerting_system_20260506) NotificationService wraps browser Notification API with localStorage permission fallback. checkHealthAndNotify compares previous vs current status per project, fires only on status transitions.
- (2026-05-07, project_tagging_workspaces_20260507) TagManager with localStorage versioning works well for client-side tag storage. getAllTags() sorted unique across repos enables filter UI. getParentDirectory() extracts grouping key from repo path.
- (2026-05-07, project_tagging_workspaces_20260507) Tag chips on project cards use click handler to toggle activeTagFilter; renderProjects() skips cards not matching filter. Directory grouping with collapsible sections using CSS toggle + hidden class works without JS state management.
- (2026-05-07, advanced_search_indexing_20260507) Search indexer initialized in main.go before API registration; BuildFromRepos blocks startup but is fast (~100ms for 22 repos). Background goroutine with time.Ticker updates index incrementally. searchFunc adapter maps internal search.SearchResult to api.SearchResult.
- (2026-05-08, cross_repo_commit_search_20260508) Avoid import cycles when adding query types: api->search->git->api cycle resolved by defining CommitSearchQuery in api package rather than search. Conversion happens in main.go adapter layer.