# Implementation Plan: Data Backup & Export

## Phase 1: Export
- [ ] Task: Build export endpoint and UI
  - [ ] Write tests for export JSON structure
  - [ ] Collect repos, tags, agents, preferences into single object
  - [ ] Add download button to settings panel
  - [ ] Include export timestamp and schema version

## Phase 2: Import
- [ ] Task: Build import endpoint and UI
  - [ ] Write tests for import validation
  - [ ] Validate schema version and required fields
  - [ ] Support merge (additive) and replace modes
  - [ ] Show confirmation with counts before applying

## Phase 3: Polish
- [ ] Task: Add error handling and manual verification
  - [ ] Write tests for malformed import rejection
  - [ ] Clear error messages for validation failures
  - [ ] Manual verification: export, clear localStorage, import, verify restored
