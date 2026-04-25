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

- [ ] **Track: Improved Search & Filtering**
  *Link: [./tracks/search-filtering_20260425/](./tracks/search-filtering_20260425/)*

- [ ] **Track: WebSocket Reliability & Security Fixes**
  *Link: [./archive/ws-reliability-fixes_20260425/](./archive/ws-reliability-fixes_20260425/)*

---

## Future Roadmap

- [ ] **Track: Enhanced Agent Orchestration & Monitoring**
  Expand the `/api/agents` capability to include real-time log streaming (via WebSockets), manual trigger overrides for cron jobs, and better error reporting for failed agent runs.

- [ ] **Track: Commit Analysis & AI Insights**
  Integrate a local or remote LLM to provide summaries of recent changes across all repositories, identifying potential bugs or architectural regressions directly in the dashboard.

- [ ] **Track: Multi-Platform Support & Advanced Git Features**
  Move beyond simple `pull` and `log` operations to support branch management, stash viewing, and multi-user authentication for shared development environments.

- [ ] **Track: Improved Search & Filtering**
  Implement server-side indexing and search across all repositories to allow finding specific code snippets or commit messages without opening individual projects.

- [ ] **Track: Unified Developer Activity Feed**
  Create a centralized view that combines Git commits, Agent logs, and potential external CI/CD status into a single timeline for better visibility of developer output.
