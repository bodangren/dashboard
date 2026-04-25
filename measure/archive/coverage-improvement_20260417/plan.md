# Implementation Plan: Agent & API Test Coverage Improvement

## Phase 1 — Assess Current Coverage

### 1.1 Run coverage baseline

- [ ] 1.1.1 Run `go test -cover ./internal/agents/...` — record baseline coverage %
- [ ] 1.1.2 Run `go test -cover ./internal/api/...` — record baseline coverage %
- [ ] 1.1.3 Run `go test -coverprofile=agents.out ./internal/agents/...` and `go tool cover -func=agents.out` to see per-function coverage
- [ ] 1.1.4 Identify untested exec-based functions: ReadCrontab, WriteCrontab, ReadLogFile, DefaultLogReader

---

## Phase 2 — Add Tests for `internal/agents`

### 2.1 Test `ReadCrontab`

- [ ] 2.1.1 Create test file `parser_test.go` if not exists
- [ ] 2.1.2 Write `TestReadCrontab_Basic` — valid crontab with agents, env vars, comments
- [ ] 2.1.3 Write `TestReadCrontab_Empty` — empty crontab returns empty slice, not nil
- [ ] 2.1.4 Write `TestReadCrontab_EnvVarsOnly` — crontab with only env vars
- [ ] 2.1.5 Write `TestReadCrontab_AgentsOnly` — crontab with only agent lines
- [ ] 2.1.6 Run `go test -cover ./internal/agents/...` — verify coverage increased

### 2.2 Test `WriteCrontab`

- [ ] 2.2.1 Create test file `writer_test.go` if not exists
- [ ] 2.2.2 Write `TestWriteCrontab_Basic` — write simple crontab, verify output format
- [ ] 2.2.3 Write `TestWriteCrontab_WithEnvVars` — write env vars, verify correct format
- [ ] 2.2.4 Write `TestWriteCrontab_RoundTrip` — read existing crontab, write it, read again, verify equal
- [ ] 2.2.5 Run `go test -cover ./internal/agents/...` — verify coverage increased

### 2.3 Test `ReadLogFile`

- [ ] 2.3.1 Write `TestReadLogFile_ValidFile` — read fake log file, verify content returned
- [ ] 2.3.2 Write `TestReadLogFile_EmptyFile` — read empty file, verify empty string returned
- [ ] 2.3.3 Write `TestReadLogFile_FileNotFound` — read non-existent file, verify error returned
- [ ] 2.3.4 Run `go test -cover ./internal/agents/...` — verify coverage

### 2.4 Test `DefaultLogReader`

- [ ] 2.4.1 Write `TestDefaultLogReader_Read` — use default reader on fake log file
- [ ] 2.4.2 Write `TestDefaultLogReader_Error` — use default reader on non-existent file
- [ ] 2.4.3 Run `go test -cover ./internal/agents/...` — verify coverage

---

## Phase 3 — Add Tests for `internal/api`

### 3.1 Identify coverage gaps in api handlers

- [ ] 3.1.1 Run `go test -coverprofile=api.out ./internal/api/...` and analyze
- [ ] 3.1.2 Focus on untested error paths in handlers

### 3.2 Add handler edge case tests

- [ ] 3.2.1 Write `TestProjectsHandler_EmptyReposDir` — repos directory exists but contains no repos
- [ ] 3.2.2 Write `TestReposHandler_EmptyResponse` — no repos discovered
- [ ] 3.2.3 Run `go test -cover ./internal/api/...` — verify coverage increased

---

## Phase 4 — Final Verification

### 4.1 Run full test suite

- [ ] 4.1.1 Run `go test -race ./...` — confirm no races
- [ ] 4.1.2 Run `go test -cover ./internal/agents/...` — verify >80%
- [ ] 4.1.3 Run `go test -cover ./internal/api/...` — verify >80%
- [ ] 4.1.4 Run `go vet ./...` — no warnings

### 4.2 Update project memory

- [ ] 4.2.1 Update `measure/tech-debt.md`: mark coverage item as resolved
- [ ] 4.2.2 Update `measure/lessons-learned.md`: add testing insights
- [ ] 4.2.3 Measure - User Manual Verification 'Phase 4.2' (Protocol in workflow.md)

### 4.3 Final commit

- [ ] 4.3.1 Commit: `test(coverage): improve agent and api package coverage to >80%`
