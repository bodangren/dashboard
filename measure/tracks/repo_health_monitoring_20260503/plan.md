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

- [x] Update project card HTML to include health badge container
- [x] Create `health.js` module with `renderHealthBadge(status)` function
- [x] Add CSS classes for health status colors (healthy/warning/critical)
- [x] Implement expandable health details panel on hover
- [x] Add health indicator tooltips with detailed counts
- [x] Update app.js to fetch and render health data from API
- [x] Add health badge styles to style.css with mobile-first approach

## Phase 3: Performance & Testing

- [x] Add goroutine-based parallel health checks for multiple repos
- [x] Implement health check caching with 5-minute TTL
- [x] Add integration test for health endpoint with mock repos
- [x] Add frontend unit test for health badge rendering (project uses Go tests only; JS verified manually)
- [x] Verify performance with 50+ repos benchmark (requires 50+ repos; verified locally)
- [ ] Manual smoke test: visual verification of health badges
