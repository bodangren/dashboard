# Implementation Plan: Agent Performance Metrics & Analytics

## Phase 1: Data Model & Storage

- [ ] Task: Extend agent run logging schema
  - [ ] Write migration tests for new `agent_runs` table
  - [ ] Create `agent_runs` table: id, agent_id, project_path, started_at, duration_ms, success, error_msg
  - [ ] Write repository tests for insert and query operations
  - [ ] Implement `AgentRunStore` with `RecordRun` and `QueryRuns` methods

## Phase 2: Aggregation Service

- [ ] Task: Build metrics calculator
  - [ ] Write tests for aggregation over various run sets
  - [ ] Implement `AgentMetrics` struct with all required fields
  - [ ] Implement `CalculateMetrics` function: total_runs, success_rate, avg_duration, p95_duration
  - [ ] Implement windowed queries (7d, 30d, all-time)
- [ ] Task: Add project-scoped metrics
  - [ ] Write tests for per-project aggregation
  - [ ] Implement `CalculateProjectMetrics` function
  - [ ] Ensure N+1 safety with batched SQL queries

## Phase 3: API Endpoints

- [ ] Task: Implement metrics handlers
  - [ ] Write handler tests with mocked store
  - [ ] Add `GET /api/agents/metrics` returning global agent metrics
  - [ ] Add `GET /api/agents/metrics?project=<path>` returning project-scoped metrics
  - [ ] Add `GET /api/agents/:id/metrics` returning single-agent history

## Phase 4: Frontend Analytics UI

- [ ] Task: Build metrics tab in agent panel
  - [ ] Write tests for metrics table rendering
  - [ ] Add "Metrics" tab next to existing agent tabs
  - [ ] Render sortable table: Agent | Runs | Success Rate | Avg Duration | Last Run
  - [ ] Add warning badges for agents exceeding thresholds
- [ ] Task: Add project-level metrics view
  - [ ] Write tests for project metric cards
  - [ ] Show per-project aggregate stats in project detail panel
  - [ ] Link to agent metrics filtered by project

## Phase 5: Integration & Verification

- [ ] Task: Wire run recording into agent execution
  - [ ] Write tests confirming run records are created on agent execution
  - [ ] Hook `RecordRun` into agent runner success/failure paths
  - [ ] Ensure failed DB writes do not break agent execution
- [ ] Task: Final verification
  - [ ] Run full test suite (`go test ./...`) — all tests must pass
  - [ ] Verify `go build` completes without errors
  - [ ] Manual browser verification of metrics tab
  - [ ] Update tech-debt.md and lessons-learned.md
  - [ ] Commit with git note
