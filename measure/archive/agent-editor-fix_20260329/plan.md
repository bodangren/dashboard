# Implementation Plan

## Phase 1: Fix Parser (TDD)

- [ ] Task 1.1: Write failing tests for opencode `-m` flag parsing
- [ ] Task 1.2: Write failing tests for `run <path>` prompt parsing
- [ ] Task 1.3: Write failing test for `>` (single) log redirect
- [ ] Task 1.4: Write failing test for full binary path extraction
- [ ] Task 1.5: Implement parser fixes to pass all new tests
- [ ] Task 1.6: Write failing test for section header capture
- [ ] Task 1.7: Implement section header capture

## Phase 2: Fix Writer

- [ ] Task 2.1: Write failing test for correct opencode command output
- [ ] Task 2.2: Write failing test for section header line in output
- [ ] Task 2.3: Implement writer fixes

## Phase 3: Fix API Handlers

- [ ] Task 3.1: Add `section_header` and `binary_path` fields to AgentJSON
- [ ] Task 3.2: Update create/update handlers to accept section header
- [ ] Task 3.3: Update tests for new fields

## Phase 4: Fix Frontend

- [ ] Task 4.1: Remove gemini/codex from harness options, hardcode opencode
- [ ] Task 4.2: Add model dropdown with searchable list (fetch from `/api/models` or use static list)
- [ ] Task 4.3: Add `/api/models` endpoint to serve model list
- [ ] Task 4.4: Display section header in agent cards
- [ ] Task 4.5: Add section header field to add/edit form

## Phase 5: Verify

- [ ] Task 5.1: Run full test suite
- [ ] Task 5.2: Run production build
- [ ] Task 5.3: Visual smoke test (manual)
