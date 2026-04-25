# Implementation Plan: Agent Project Sections — Crontab Reorganization

## Phase 1 — Writer: ReorganizeAutomation

- [x] 1.1 Add separator constant and `findSeparator` / `ensureSeparator` helpers to `writer.go`
- [x] 1.2 Implement `Crontab.ReorganizeAutomation(projects []string)`: split Lines into preamble (above separator) and automation (below), collect existing agents by directory, rebuild automation section with headers for all projects, preserve agent order within each project
- [x] 1.3 Write `TestReorganizeAutomation`: given a messy crontab and a project list, verify separator exists, all projects have headers, agents are grouped, preamble is untouched

## Phase 2 — Writer: Fix AddAgent / DeleteAgent / UpdateAgent

- [x] 2.1 Modify `AddAgent`: call `ensureSeparator`, append agent to Lines, call `ReorganizeAutomation` to place under correct section header
- [x] 2.2 Modify `DeleteAgent`: remove agent line; remove section header only if next line is not another agent in the same section
- [x] 2.3 Modify `UpdateAgent`: if directory changed, delete from old section and re-add under new section via `ReorganizeAutomation`
- [x] 2.4 Write `TestAddAgentUnderExistingSection`: add two agents for same project, verify single header
- [x] 2.5 Write `TestDeleteAgentPreservesHeader`: add two agents, delete one, verify header remains
- [x] 2.6 Write `TestDeleteLastAgentRemovesHeader`: add one agent, delete it, verify header removed
- [x] 2.7 Write `TestUpdateAgentMovesProject`: change agent's directory, verify it moves to new section

## Phase 3 — Handler: Pass Projects to AgentHandler

- [x] 3.1 Add `repos []string` field to `AgentHandler` and `SetRepos([]string)` method
- [x] 3.2 In `main.go`, call `agentHandler.SetRepos(repos)` after scanning
- [x] 3.3 Modify `createAgent`, `deleteAgent`, `updateAgent` handlers to call `ct.ReorganizeAutomation(ah.repos)` before writing crontab
- [x] 3.4 Write `TestAgentCreateHandlerUsesRepos` verifying all repos get section headers

## Phase 4 — Frontend: Pass Projects on Create/Edit

- [x] 4.1 Not needed — backend already receives repos from `main.go` via `SetRepos`
- [x] 4.2 N/A
- [x] 4.3 N/A

## Phase 5 — Cleanup & Verify

- [x] 5.1 Run all tests: `go test ./...` — all pass
- [x] 5.2 Build binary: `go build -o dashboard .` — succeeds
- [x] 5.3 Delete empty stub `measure/tracks/agent-project-sections_20260330/`
