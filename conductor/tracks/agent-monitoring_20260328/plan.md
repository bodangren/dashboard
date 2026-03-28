# Implementation Plan: Autonomous Agent Monitoring

## Phase 1 — Backend: Crontab Parser & Agent Types

- [x] 1.1 Define `Agent` struct in `internal/agents/types.go`: Schedule, Directory, Harness, Model, Prompt, LogPath, Enabled, LineIndex (position in crontab for round-tripping)
- [x] 1.2 Write test: `TestParseCrontab` — given sample crontab output, parser returns correct agent structs with enabled/disabled state
- [x] 1.3 Implement `internal/agents/parser.go`: read crontab output, classify lines (env/var, comment, agent, other), extract agent fields using regex for each harness pattern
- [x] 1.4 Write test: `TestParseCrontabPreservesNonAgent` — verify non-agent lines (env vars, update jobs, blank lines) are preserved in the parsed structure
- [x] 1.5 Implement round-trip data structure: store full crontab as slice of lines with metadata, agents reference their line indices
- [x] 1.6 Run tests: `go test ./internal/agents/...`

## Phase 2 — Backend: Crontab Read/Write

- [x] 2.1 Write test: `TestWriteCrontab` — modifying an agent and writing back preserves all non-agent lines
- [x] 2.2 Implement `internal/agents/writer.go`: serialize agents back to crontab lines, merge with preserved non-agent lines, write via `crontab -`
- [x] 2.3 Write test: `TestToggleAgent` — enable/disable toggles comment prefix on the agent's crontab line
- [x] 2.4 Implement toggle: comment/uncomment an agent's line while preserving the rest
- [x] 2.5 Write test: `TestAddAgent` — new agent appended correctly with valid cron syntax
- [x] 2.6 Implement add: construct crontab line from Agent struct, append to crontab
- [x] 2.7 Write test: `TestDeleteAgent` — agent line removed, other lines intact
- [x] 2.8 Implement delete: remove agent's line from the crontab
- [x] 2.9 Run tests: `go test ./internal/agents/... -cover`

## Phase 3 — Backend: Agent API Handlers

- [ ] 3.1 Write test: `TestAgentsHandler` — mock crontab reader, verify JSON response shape
- [ ] 3.2 Implement `GET /api/agents` — read crontab, parse agents, return JSON array grouped by project
- [ ] 3.3 Write test: `TestAgentCreateHandler` — valid POST creates agent in crontab
- [ ] 3.4 Implement `POST /api/agents` — parse JSON body, add agent, write crontab
- [ ] 3.5 Write test: `TestAgentUpdateHandler` — PUT updates agent in-place
- [ ] 3.6 Implement `PUT /api/agents/{index}` — update agent fields, rewrite crontab
- [ ] 3.7 Write test: `TestAgentDeleteHandler` — DELETE removes agent
- [ ] 3.8 Implement `DELETE /api/agents/{index}` — remove agent, rewrite crontab
- [ ] 3.9 Write test: `TestAgentToggleHandler` — PATCH toggles enabled/disabled
- [ ] 3.10 Implement `PATCH /api/agents/{index}/toggle` — flip comment state
- [ ] 3.11 Write test: `TestAgentLogHandler` — returns last N lines of agent log file
- [ ] 3.12 Implement `GET /api/agents/{index}/log` — read log file, return last N lines
- [ ] 3.13 Register new routes in `main.go`
- [ ] 3.14 Run tests: `go test ./internal/api/... -cover`

## Phase 4 — Frontend: Navigation & Agent List Page

- [ ] 4.1 Add nav bar to `index.html`: "Commits" and "Agents" tabs in the header
- [ ] 4.2 Add nav bar to `agents.html`: same nav, "Agents" tab active
- [ ] 4.3 Add nav styles to `style.css`: tab bar styling consistent with dark terminal theme
- [ ] 4.4 Create `static/agents.html`: page shell with agent list container
- [ ] 4.5 Create `static/agents.js`: fetch `/api/agents`, render grouped agent cards
- [ ] 4.6 Each agent card displays: project name, harness, model, schedule, enabled badge, log path
- [ ] 4.7 Disabled agents styled distinctly (dimmed, strikethrough schedule)

## Phase 5 — Frontend: CRUD Forms & Schedule Editor

- [ ] 5.1 Implement "Add Agent" form (modal or inline): directory select, harness select, model input, schedule editor, prompt input, log path input
- [ ] 5.2 Populate directory select from `/api/projects` repo list
- [ ] 5.3 Implement schedule editor: human-friendly presets ("Every N hours", "Daily at", "Custom cron") with live preview
- [ ] 5.4 Implement edit mode: click agent card to open pre-filled edit form
- [ ] 5.5 Implement delete: delete button with confirmation dialog
- [ ] 5.6 Implement toggle: switch/button on each agent card to enable/disable

## Phase 6 — Frontend: Log Viewer & Polish

- [ ] 6.1 Implement expandable log section per agent: fetch `/api/agents/{index}/log`, display last 50 lines in a scrollable pre block
- [ ] 6.2 Show last run time (from log file mtime) on each agent card
- [ ] 6.3 Handle missing log files gracefully ("No log file found")
- [ ] 6.4 Visual polish: consistent with terminal theme, compact layout
- [ ] 6.5 Manual test: full CRUD flow — add, edit, toggle, delete agents via UI

## Phase 7 — Integration & Final Verification

- [ ] 7.1 Write integration test: `TestAgentCRUDEndToEnd` — start server, exercise all agent API endpoints against real crontab mock
- [ ] 7.2 Run full test suite: `go test ./... -cover` — confirm ≥ 80% coverage on new packages
- [ ] 7.3 Run `go vet ./...` — no issues
- [ ] 7.4 Manual test: verify against actual crontab, confirm non-agent lines preserved
- [ ] 7.5 Commit: `feat(agents): autonomous agent monitoring with crontab CRUD`
