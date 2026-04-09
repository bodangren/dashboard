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
| 2026-04-09 | critical-bugs-rewrite_20260406 | ARCH-01: Duplicate api/git.Commit types still exist | Low | Open | Helper function in main.go works around manual mapping; full merge would require shared types package |
| 2026-04-09 | critical-bugs-rewrite_20260406 | ARCH-02: /api/repos lightweight endpoint added | Medium | Resolved | Frontend agents.js still uses /api/projects — needs update in Phase 3 |
