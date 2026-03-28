# Lessons Learned

> This file is curated working memory, not an append-only log. Keep it at or below **50 lines**.
> Remove or condense entries that are no longer relevant to near-term planning.

## Architecture & Design

- (2026-03-28, git-view-enhance_20260328) CSS grid `auto-fill` with small minmax was too cramped at 3 cols for dense commit info. Fixed 2-col grid is better for terminal-styled dashboards with monospace text.
- (2026-03-29, agent-editor-fix_20260329) Section header comments in crontab are separate `Line` entries in the slice, not inline with the agent. When adding new agents with section headers, insert a `LineComment` before the `LineAgent`. When deleting, remove the preceding comment too.

## Recurring Gotchas

- (2026-03-28, git-view-enhance_20260328) Pure CSS/JS tracks don't have unit test coverage in a Go project. Manual verification is the only gate — always defer manual-smoke-test tasks until user can visually confirm.
- (2026-03-29, agent-editor-fix_20260329) Real crontab uses `>` (single) for redirect, not `>>`. Regex must handle both. OpenCode uses `-m` flag and `run <path>` positional, not `--model`/`--prompt`.

## Patterns That Worked Well

- (2026-03-28, git-view-enhance_20260328) Keeping `.commit-age-badge` as a separate DOM element (not innerHTML string concat) made it easy to conditionally append only when commits exist.
- (2026-03-29, agent-editor-fix_20260329) Pending-comment pattern in parser — track the last comment seen, attach it to the next agent line, reset on non-comment/non-agent lines. Clean way to capture section headers without modifying the Line struct heavily.

## Planning Improvements

- (2026-03-28, git-view-enhance_20260328) Responsive column count is subjective — spec said 3→2→1 but user found 3 too cramped. Should prototype layout decisions before committing to column counts in the spec.
