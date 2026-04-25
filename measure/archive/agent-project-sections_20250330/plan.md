# Implementation Plan: Fix Agent Project Sections

## Phase 1 — Revert Commits View Changes

- [x] 1.1 Revert internal/api/handlers.go to remove Agents field from Project struct
- [x] 1.2 Revert static/app.js to remove agent indicator rendering
- [x] 1.3 Revert static/style.css to remove agent indicator styles

## Phase 2 — Implement Project Sections in Agents View

- [x] 2.1 Modify loadAgents() to fetch both projects and agents
- [x] 2.2 Group agents by project directory
- [x] 2.3 Render all projects as sections (even empty ones)
- [x] 2.4 Add attachProjectSelectHandler() to auto-update section header

## Phase 3 — Fix Frontend-Backend Field Mapping

- [x] 3.1 Change logPath to log_path in agents.js request body

## Phase 4 — Reformat Crontab

- [x] 4.1 Add separator line # ==AUTOMATION BELOW THIS LINE==
- [x] 4.2 Reorganize agents under project section headers
- [x] 4.3 Install updated crontab

## Phase 5 — Verify

- [x] 5.1 Run all tests
- [x] 5.2 Build dashboard binary
- [x] 5.3 Verify agents view shows project sections
- [x] 5.4 Verify commits view is restored
