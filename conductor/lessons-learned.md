# Lessons Learned

> This file is curated working memory, not an append-only log. Keep it at or below **50 lines**.
> Remove or condense entries that are no longer relevant to near-term planning.

## Architecture & Design
<!-- Decisions made that future tracks should be aware of -->

- (2026-03-28, git-view-enhance_20260328) CSS grid `auto-fill` with small minmax was too cramped at 3 cols for dense commit info. Fixed 2-col grid is better for terminal-styled dashboards with monospace text.

## Recurring Gotchas
<!-- Problems encountered repeatedly; save future tracks from the same pain -->

- (2026-03-28, git-view-enhance_20260328) Pure CSS/JS tracks don't have unit test coverage in a Go project. Manual verification is the only gate — always defer manual-smoke-test tasks until user can visually confirm.

## Patterns That Worked Well
<!-- Approaches worth repeating -->

- (2026-03-28, git-view-enhance_20260328) Keeping `.commit-age-badge` as a separate DOM element (not innerHTML string concat) made it easy to conditionally append only when commits exist.

## Planning Improvements
<!-- Notes on where estimates were wrong and why -->

- (2026-03-28, git-view-enhance_20260328) Responsive column count is subjective — spec said 3→2→1 but user found 3 too cramped. Should prototype layout decisions before committing to column counts in the spec.
