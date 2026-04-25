# Specification: Agent & API Test Coverage Improvement

## Track ID
`coverage-improvement_20260417`

## Priority
**P1 — Important. Improves reliability and enables confident refactoring.**

## Overview

Test coverage for `internal/agents` and `internal/api` packages is below 80%. The exec-based functions (`ReadCrontab`, `WriteCrontab`, `ReadLogFile`, `DefaultLogReader`) are completely untested. This track adds integration tests using mock crontabs and fake git directories to close the coverage gap.

---

## Coverage Gaps

### COV-01: `ReadCrontab` is untested
- **File:** `internal/agents/parser.go`
- **Current coverage:** 0%
- **Problem:** `ReadCrontab` parses crontab output but has no unit tests
- **Fix:** Add tests with various crontab formats (comments, env vars, agents, empty sections)

### COV-02: `WriteCrontab` is untested  
- **File:** `internal/agents/writer.go`
- **Current coverage:** 0%
- **Problem:** `WriteCrontab` writes crontab entries but has no unit tests
- **Fix:** Add tests with mock crontab content, verify correct output format

### COV-03: `ReadLogFile` is untested
- **File:** `internal/agents/parser.go`
- **Current coverage:** 0%
- **Problem:** Log file reading for agent output has no tests
- **Fix:** Add tests with fake log files, various content sizes

### COV-04: `DefaultLogReader` is untested
- **File:** `internal/agents/parser.go`
- **Current coverage:** 0%
- **Problem:** Default log reader implementation has no tests
- **Fix:** Add tests with mock log content

### COV-05: `GetCommitInfo` edge cases untested
- **File:** `internal/git/git.go`
- **Current coverage:** Partial
- **Problem:** `GetCommitInfo` handles empty hashes, invalid paths but not fully tested
- **Fix:** Add tests for error cases

---

## Acceptance Criteria

1. `go test -cover ./internal/agents/...` shows >80% coverage
2. `go test -cover ./internal/api/...` shows >80% coverage  
3. All exec-based functions have at least one integration test
4. `go test -race ./...` passes with no races
5. No existing tests are broken by the changes

---

## Out of Scope

- Adding unit tests for already-well-tested functions
- Changing production code behavior
- Adding new public APIs
- Performance benchmarking
