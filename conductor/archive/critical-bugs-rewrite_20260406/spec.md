# Specification: Critical Bugs & Display Rewrite

## Track ID
`critical-bugs-rewrite_20260406`

## Priority
**P0 — Highest Priority. Blocks all future feature work.**

## Overview

A comprehensive audit of the codebase revealed critical bugs, data races, fragile parsing logic, and significant display/UX issues across both backend and frontend. This track addresses all identified problems in dependency order: backend correctness first, then data layer cleanup, then frontend display overhaul.

---

## Bugs (Must Fix)

### BUG-01: Race condition on package-level function variables
- **File:** `internal/api/handlers.go:127-134, 177-185`
- **Problem:** `defaultGetCommits`, `defaultGetDiff`, `defaultPullRepo` are package-level `var`s mutated by `SetGitFuncs()`/`SetPullFunc()`. The `Handler` struct captures references to these at creation via `RegisterRoutes`, but concurrent HTTP requests can read stale or nil function pointers if the globals are mutated during request processing. This is a data race under Go's race detector.
- **Fix:** Eliminate global mutable state. Pass function implementations directly into `RegisterRoutes` or via a `HandlerConfig` struct. The handler should store function fields that are set once at construction and never mutated.

### BUG-02: `detectHarness` uses fragile regex-string slicing
- **File:** `internal/agents/parser.go:151-158`
- **Problem:** `p.String()[2 : len(p.String())-2]` strips `\b` delimiters from the regex source string. This only works because all patterns happen to be `\b<word>\b`. Any regex change (e.g. `^opencode`, `(?:opencode)`) silently breaks harness detection, returning an empty string and producing malformed crontab output.
- **Fix:** Replace with an explicit name map from compiled regex to harness string, or use named capture groups. The harness name must not be derived from the regex source string.

### BUG-03: `isEnvVarLine` has wrong operator precedence
- **File:** `internal/agents/parser.go:106-108`
- **Problem:** `strings.Contains(line, "=") && !strings.Contains(line, " ") || regexp.MustCompile(...)` evaluates as `(A && B) || C` due to Go's operator precedence (`||` and `&&` are left-to-right but `&&` binds tighter — actually this IS correct in Go, but the expression is misleading and the first clause alone is wrong). The real problem: `strings.Contains(line, "=") && !strings.Contains(line, " ")` will match `PATH=/usr/local/bin` but also any token like `key=` with no spaces, potentially misclassifying cron command arguments containing `=` (e.g. `make test=1`).
- **Fix:** Rely solely on the regex `^[A-Za-z_][A-Za-z0-9_]*=`. Remove the first clause entirely or parenthesize correctly. Add a test with a cron line containing `=` in the command portion.

### BUG-04: `toggleAgent` fragile agent re-lookup
- **File:** `internal/api/agent_handlers.go:236-248`
- **Problem:** After toggling, the handler re-parses the crontab and tries to find the toggled agent with: `if a.LineIndex == agentList[idx].LineIndex || a == agentList[idx)`. The `||` makes the `LineIndex` check redundant (pointer equality is already checked). If the toggle changes line count (it doesn't currently, but the code doesn't guarantee this), `idx` could be out of bounds, leading to a nil pointer dereference at line 243: `agent = agentList[idx]`.
- **Fix:** Simplify: after toggle and write, re-read the crontab, find the agent by its original line content or a stable identifier, and return the updated state. Remove the confusing `||` logic.

### BUG-05: Agent TOCTOU — index-based agent referencing
- **File:** `internal/api/agent_handlers.go` (all mutation handlers), `static/agents.js`
- **Problem:** The frontend identifies agents by their array index (`/api/agents/0`, `/api/agents/1`). If the crontab is modified between listing and acting (e.g. two browser tabs, or external `crontab -e`), the index may point to a different agent or be out of bounds. This can cause the wrong agent to be deleted/toggled/updated, or corrupt the crontab entirely.
- **Fix:** Add a stable `id` field to agents derived from a hash of their immutable properties (schedule + directory + model). The frontend sends this id instead of an array index. The backend resolves the id to the correct agent on each mutation.

---

## Architecture / Data Layer Issues

### ARCH-01: Duplicate `api.Commit` vs `git.Commit` type
- **File:** `main.go:37-53`, `internal/api/handlers.go:13-20`
- **Problem:** `api.Commit` and `git.Commit` are identical in shape. `main.go` manually maps every field in a loop. This is boilerplate that must be updated in lockstep whenever fields change.
- **Fix:** Have the `api` package accept the `git.Commit` type directly (introduce a shared interface or have `api` import from a shared types package). Alternatively, use a simple conversion function in the `git` package.

### ARCH-02: Redundant `/api/projects` call from agents page
- **File:** `static/agents.js:112-113`
- **Problem:** The agents page calls `fetch('/api/projects')` to get the project list for grouping. This endpoint runs `git log` on every repo (up to 10 commits each). For the agents page, only repo paths are needed — the commit history is wasted work.
- **Fix:** Add a lightweight `/api/repos` endpoint that returns only `{name, path}` without running git commands. Update `agents.js` to use this endpoint.

### ARCH-03: Empty repos are invisible
- **File:** `internal/api/handlers.go:80`
- **Problem:** `if err != nil || len(commits) == 0 { continue }` silently skips repos with no commits. New or freshly cloned repos don't appear in the dashboard, contradicting the product spec's "auto-discovery" feature.
- **Fix:** Always include discovered repos in the response. For repos with no commits, return an empty commits array. Show them in the UI with a "no commits yet" indicator.

### ARCH-04: Service Worker caches API responses
- **File:** `static/sw.js:53-69`
- **Problem:** The service worker caches `/api/` responses with a network-first strategy. Stale data can be served from cache when the network response fails. For a real-time dashboard, this is counterproductive — users see old commit data.
- **Fix:** Remove API response caching from the service worker. Only cache static assets (HTML, CSS, JS).

### ARCH-05: `AddAgent` calls `ReorganizeAutomation(nil)` standalone
- **File:** `internal/agents/writer.go:164`
- **Problem:** When `AddAgent` is called outside of the handler (e.g. in tests or future code), `ReorganizeAutomation(nil)` means no project list is available for section grouping.
- **Fix:** Make `AddAgent` accept an optional project list, or document that `ReorganizeAutomation` should be called separately with the project list.

---

## Display / Frontend Issues

### UX-01: Full filesystem paths shown in project cards
- **File:** `static/app.js:44-45`
- **Problem:** Project cards display `/home/daniel-bo/Desktop/dashboard` as the path. This is noise — the project name is already shown. The full path adds visual clutter and leaks filesystem structure.
- **Fix:** Show only the project name prominently. Display the path only on hover/tooltip or as a truncated relative path (e.g. `~/Desktop/dashboard`).

### UX-02: Hardcoded binary path in agent form
- **File:** `static/agents.js:316`
- **Problem:** `/home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode` is hardcoded as the default binary path. This breaks on any other machine and will break when Node is upgraded.
- **Fix:** Default to just `opencode` (rely on `$PATH`). If the existing agent has a binary path, preserve it. Remove the hardcoded absolute path.

### UX-03: Agent timing visualization is opaque
- **File:** `static/agents.js:50-66`
- **Problem:** 7 day blocks + 24 hour blocks + minute number is visually dense with no labels. New users cannot tell what the colored blocks represent without prior knowledge.
- **Fix:** Add a small legend/header: "Days: S M T W T F S" above day blocks and "Hours" label. Make the minute display more prominent. Consider showing the human-readable schedule text prominently and the visualization as secondary.

### UX-04: Diff page lacks commit message context
- **File:** `static/diff.html:19`, `static/diff.js:70-78`
- **Problem:** The diff page title shows only the hash (e.g. "abc1234"). The commit message is not displayed at all because the `/api/diff` response doesn't include it. The meta section checks for `data.author`/`data.message`/`data.timestamp` which are never returned by `DiffResponse`.
- **Fix:** Extend `DiffResponse` to include `message`, `author`, and `timestamp` fields. Populate them from a `git log -1` call alongside the diff. Display them prominently on the diff page.

### UX-05: Duplicate `esc()` function across three files
- **File:** `static/app.js:110-116`, `static/agents.js:6-12`, `static/diff.js:20-26`
- **Problem:** The HTML escape function is copy-pasted identically in three files. Any bug fix or enhancement must be applied three times.
- **Fix:** Create `static/utils.js` with shared utility functions (`esc`, `relativeTime`, etc.). Include it via `<script>` tag before the page-specific scripts.

### UX-06: `relativeTime` doesn't handle future timestamps
- **File:** `static/app.js:20-30`
- **Problem:** If a commit has a future timestamp (clock skew, timezone mismatch), `Date.now() - new Date(isoStr).getTime()` is negative. This produces "NaNs ago" or "-5s ago".
- **Fix:** Clamp the diff to minimum 0, or show "just now" for future timestamps.

### UX-07: CSS media query sprawl — ~50% of rules are mobile overrides
- **File:** `static/style.css`
- **Problem:** Almost every CSS rule has an identical `@media (max-width: 768px)` block immediately after it. This doubles the CSS size and makes maintenance hard.
- **Fix:** Refactor to mobile-first CSS (base styles = mobile, `@media (min-width: 768px)` for desktop). Use CSS custom properties for responsive spacing. Consolidate related breakpoints.

---

## Acceptance Criteria

1. All P0 bugs (BUG-01 through BUG-05) are fixed with regression tests
2. `go test -race ./...` passes with no data races
3. All existing tests continue to pass
4. Empty repos appear in the dashboard
5. Agents page loads without triggering `git log` on every repo
6. Agent mutations use stable identifiers, not array indices
7. Diff page shows commit message, author, and date
8. Project cards show clean names without full filesystem paths
9. Service worker no longer caches API responses
10. No duplicate utility functions across frontend files
11. `detectHarness` does not derive harness name from regex source string
12. All new code has >80% test coverage

---

## Out of Scope

- Real-time log streaming (on roadmap)
- LLM-powered commit summaries (on roadmap)
- Multi-user authentication
- Push/write operations to repositories
- Pagination of commit history
- WebSocket support
- Any new features beyond fixing identified issues
