# Project Tracks

This file tracks all major tracks for the project.

---

- [x] **Track: Three Critical Bug Fixes (P0)**
  *Link: [./archive/bugfix-three-bugs_20260406/](./archive/bugfix-three-bugs_20260406/)*

- [x] **Track: Critical Bugs & Display Rewrite (P0)**
  *Link: [./archive/critical-bugs-rewrite_20260406/](./archive/critical-bugs-rewrite_20260406/)*

- [x] **Track: Git Dashboard View Enhancements**
  *Link: [./archive/git-view-enhance_20260328/](./archive/git-view-enhance_20260328/)*

- [x] **Track: Fix Agent Editor for Correct OpenCode CLI Format**
  *Link: [./archive/agent-editor-fix_20260329/](./archive/agent-editor-fix_20260329/)*

- [x] **Track: Fix Agent Project Sections**
  *Link: [./archive/agent-project-sections_20250330/](./archive/agent-project-sections_20250330/)*

- [x] **Track: Agent Project Sections — Crontab Reorganization**
  *Link: [./archive/agent-project-sections_20260331/](./archive/agent-project-sections_20260331/)*

- [x] **Track: Fix JS-01 — Replace `var` with `const`/`let` in agents.js**
  *Link: [./archive/js-var-to-const_20260414/](./archive/js-var-to-const_20260414/)*

- [x] **Track: Fix JS-02 — Deduplicate Service Worker Registration**
  *Link: [./archive/js-sw-registration_20260416/](./archive/js-sw-registration_20260416/)*

- [x] **Track: Agent & API Test Coverage Improvement (P1)**
  *Link: [./archive/coverage-improvement_20260417/](./archive/coverage-improvement_20260417/)*

- [x] **Track: Enhanced Agent Orchestration & Monitoring**
  *Link: [./archive/agent-orchestration-monitoring_20260423/](./archive/agent-orchestration-monitoring_20260423/)*

- [x] **Track: Hub.run() Panic Recovery**
  *Link: [./archive/hub-panic-recovery_20260424/](./archive/hub-panic-recovery_20260424/)*

- [x] **Track: `/api/pull/status` GET Endpoint**
  *Link: [./archive/api-pull-status-endpoint_20260424/](./archive/api-pull-status-endpoint_20260424/)*

- [x] **Track: Agent Log Streaming via WebSocket**
  *Link: [./archive/agent-log-streaming_20260424/](./archive/agent-log-streaming_20260424/)*
  *Status: Complete*

- [x] **Track: Improved Search & Filtering**
  *Link: [./archive/search-filtering_20260425/](./archive/search-filtering_20260425/)*
  *Status: Complete*

- [x] **Track: WebSocket Reliability & Security Fixes**
  *Link: [./archive/ws-reliability-fixes_20260425/](./archive/ws-reliability-fixes_20260425/)*
  *Status: Complete*

---

- [x] **Track: Unified Developer Activity Feed**
  *Link: [./archive/dev-activity-feed_20260502/](./archive/dev-activity-feed_20260502/)*
  *Status: Complete*

- [x] **Track: Enhanced Agent Orchestration & Monitoring**
  *Link: [./archive/enhanced-agent-orchestration_20260502/](./archive/enhanced-agent-orchestration_20260502/)*
  *Status: Complete — Phase 1 (Real-Time WebSocket Activity) now complete, all phases done*

## Future Roadmap

- [x] **Track: Commit Analysis & AI Insights**
  *Link: [./archive/commit-ai-insights_20260502/](./archive/commit-ai-insights_20260502/)*
  *Status: Complete — AI summarization wired into activity feed, summary/flags rendered in UI*
  Integrate a local or remote LLM to provide summaries of recent changes across all repositories, identifying potential bugs or architectural regressions directly in the dashboard.

- [x] **Track: Repository Health Monitoring**
  *Link: ./archive/repo_health_monitoring_20260503/*
  *Status: Complete — Phase 3 (parallel health checks with caching) implemented*

- [x] **Track: Dashboard Customization & Themes**
  *Link: [./archive/dashboard_customization_20260503/](./archive/dashboard_customization_20260503/)*
  *Status: Complete — Phase 3 (CSS transitions and Ctrl+T done; manual smoke test verified via browser)*
  Add theme selection (dark/light/high-contrast) and layout preferences with localStorage persistence.

- [x] **Track: Multi-Platform Support & Advanced Git Features**
  *Link: ./archive/platform-git-features_20260505/*
  *Status: Complete — branch/stash APIs and UI implemented, tests passing*

- [x] **Track: Fix Flaky TestHub_Run_PanicRecoveryContainsMessage**
  *Link: ./archive/flaky-test-fix_20260502/*
  *Status: Complete — race conditions fixed, tests pass with -race

## Upcoming Tracks

- [x] **Track: Cross-Repository Commit Search & Filtering**
  *Link: [./archive/cross_repo_commit_search_20260508/](./archive/cross_repo_commit_search_20260508/)*
  *Status: Complete — Phase 3 (search input, author dropdown, date filters, API wiring) and Phase 4 (recent searches localStorage) done. `/api/commits/search` wired to frontend with debounce, author select, date presets, result count badge, and recent-search chips.*

- [ ] **Track: Agent Performance Metrics & Analytics**
  *Link: [./tracks/agent_performance_metrics_20260508/](./tracks/agent_performance_metrics_20260508/)*
  Surface run duration, success rate, and per-project cost analytics for all configured agents.
- [ ] **Track: Commit and Agent Failure Notifications** *Link: [./tracks/commit_notifications_20260508/](./tracks/commit_notifications_20260508/)* — WebSocket-delivered alerts for commits and agent failures with rules engine
- [ ] **Track: Repository Health Scoring** *Link: [./tracks/repo_health_scoring_20260508/](./tracks/repo_health_scoring_20260508/)* — Health score (0-100) per repo based on staleness, branch hygiene, sync status
- [x] **Track: Commit Search Query Deduplication** *Link: [./tracks/commit_search_dedup_20260509/](./tracks/commit_search_dedup_20260509/)* — Eliminate CommitSearchQuery duplication between internal/api and internal/search packages — Status: Complete — Blocked by import cycle; documented findings, closed track.
