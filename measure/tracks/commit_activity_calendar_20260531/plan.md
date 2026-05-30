# Plan: Commit Activity Calendar and Streak Tracking

## Phase 1: Data Layer — Commit History Aggregation

### Tests First
- [ ] Test `GetDailyCommitCounts` returns map[date]count for single repo
- [ ] Test `GetDailyCommitCounts` aggregates across multiple repos
- [ ] Test `GetDailyCommitCounts` handles empty repos (no panic)
- [ ] Test `GetDailyCommitCounts` deduplicates commits by hash across repos
- [ ] Test `GetDailyCommitCounts` respects 365-day window

### Implementation
- [ ] Add `GetDailyCommitCounts(repos []string, days int) map[string]int` in `git` package
- [ ] Use `git log --format=%ad --date=short --all` per repo
- [ ] Aggregate into date→count map
- [ ] Add `/api/activity/calendar?days=365` endpoint

## Phase 2: Streak Calculation

### Tests First
- [ ] Test `CalculateStreaks` with continuous daily commits
- [ ] Test `CalculateStreaks` with gaps (resets current streak)
- [ ] Test `CalculateStreaks` with single-day streak
- [ ] Test `CalculateStreaks` with no commits (zero streaks)

### Implementation
- [ ] Add `CalculateStreaks(counts map[string]int) (current int, longest int)` in `git` or `api` package
- [ ] Include streak data in `/api/activity/calendar` response

## Phase 3: Frontend Calendar Component

### Tests (Manual QA)
- [ ] Calendar renders 52 weeks × 7 days
- [ ] Color intensity matches commit count tiers (0, 1-2, 3-5, 6-10, 10+)
- [ ] Current streak badge displays correctly
- [ ] Longest streak badge displays correctly
- [ ] Mobile: horizontal scroll or compact view

### Implementation
- [ ] Create `calendar.js` module
- [ ] Generate SVG or CSS grid calendar
- [ ] Add color scale CSS custom properties
- [ ] Wire click handler to filter commits by date
- [ ] Add calendar section to `index.html`

## Phase 4: Integration and Polish

- [ ] Add calendar section toggle in preferences
- [ ] Update `relativeTime` to handle future dates gracefully
- [ ] Run `go test ./...`
- [ ] Manual verification in browser
- [ ] Update `measure/tech-debt.md` and `measure/lessons-learned.md`
- [ ] Commit and push
