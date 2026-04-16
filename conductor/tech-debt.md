# Tech Debt Registry

> This file is curated working memory, not an append-only log. Keep it at or below **50 lines**.
> Remove or summarize entries that are no longer relevant to near-term planning.
>
> **Severity:** `Critical` | `High` | `Medium` | `Low`
> **Status:** `Open` | `Resolved`

| Date | Track | Item | Severity | Status | Notes |
|------|-------|------|----------|--------|-------|
| 2026-01-01 | example_track | Example: Hardcoded timeout value | Low | Resolved | Replaced with config value in v1.2 |
| 2026-03-28 | agent-monitoring | agents & api coverage below 80% — exec-based funcs (ReadCrontab, WriteCrontab, ReadLogFile, DefaultLogReader) are untested | Low | Open | Integration test with real crontab mock could close the gap |
| 2026-04-06 | bugfix-three-bugs | Phase 2.2 deferred: `/api/pull/status` GET endpoint not implemented | Low | Open | Added to future roadmap in critical-bugs-rewrite track |
| 2026-04-09 | critical-bugs-rewrite_20260406 | BUG-01 through BUG-05: Race conditions, fragile regex, index-based agent referencing | Critical | Resolved | HandlerConfig, explicit harness map, agent IDs implemented |
| 2026-04-10 | critical-bugs-rewrite_20260406 | ARCH-01: Duplicate api/git.Commit types | Low | Resolved | ToAPICommit() method on git.Commit eliminates manual mapping |
| 2026-04-10 | critical-bugs-rewrite_20260406 | ARCH-02: /api/repos lightweight endpoint | Medium | Resolved | Frontend agents.js now uses /api/repos |
| 2026-04-10 | critical-bugs-rewrite_20260406 | ARCH-03: Empty repos invisible | Medium | Resolved | Handler now includes repos with zero commits |
| 2026-04-10 | critical-bugs-rewrite_20260406 | ARCH-04: Service worker cached API | Medium | Resolved | Removed API caching from service worker |
| 2026-04-10 | critical-bugs-rewrite_20260406 | ARCH-05: AddAgent ReorganizeAutomation(nil) | Low | Resolved | Now passes empty slice for consistent behavior |
| 2026-04-11 | critical-bugs-rewrite_20260406 | ARCH-06: GetCommitInfo for diff metadata | Low | Resolved | git.GetCommitInfo returns message/author/timestamp, wired into DiffResponse |
| 2026-04-12 | critical-bugs-rewrite_20260406 | UX-07: Mobile-first CSS refactor | Low | Resolved | Full conversion done: mobile base + desktop @media (min-width: 769px) pattern, CSS custom properties, consolidated queries |
| 2026-04-14 | js-var-to-const_20260414 | JS-01: agents.js uses `var` throughout (styleguide forbids it) | Medium | Resolved | 93 var → 84 const + 9 let (loop counters i/h/d, cachedModels/Repos, hour/dayCheckboxes). Zero var remaining. |
| 2026-04-16 | js-sw-registration_20260416 | JS-02: Service worker registration duplicated in app.js and diff.js | Low | Resolved | Extracted to registerServiceWorker() in utils.js, called from app.js and diff.js |
