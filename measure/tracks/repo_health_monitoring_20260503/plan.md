# Repository Health Monitoring Plan

## Phase 1: Backend Health Indicators

- [x] Create `health.go` package with health indicator functions
- [x] Implement `GetDirtyStatus(repoPath) DirtyStatus` - count staged/unstaged/untracked files
- [x] Implement `GetBranchDivergence(repoPath) BranchDivergence` - ahead/behind remote counts
- [x] Implement `GetStaleBranches(repoPath, threshold) StaleBranchInfo` - branches older than threshold
- [x] Add unit tests for each health indicator function with temp git repos
- [x] Create `HealthResponse` struct in types.go with all health fields
- [x] Add `/api/health?repo=<name>` endpoint to handler
- [x] Integrate health checks into `/api/projects` response for all repos

## Phase 2: Frontend Health Display

- [ ] Update project card HTML to include health badge container
- [ ] Create `health.js` module with `renderHealthBadge(status)` function
- [ ] Add CSS classes for health status colors (healthy/warning/critical)
- [ ] Implement expandable health details panel on hover
- [ ] Add health indicator tooltips with detailed counts
- [ ] Update app.js to fetch and render health data from API
- [ ] Add health badge styles to style.css with mobile-first approach

## Phase 3: Performance & Testing

- [ ] Add goroutine-based parallel health checks for multiple repos
- [ ] Implement health check caching with 5-minute TTL
- [ ] Add integration test for health endpoint with mock repos
- [ ] Add frontend unit test for health badge rendering
- [ ] Verify performance with 50+ repos benchmark
- [ ] Manual smoke test: visual verification of health badges
