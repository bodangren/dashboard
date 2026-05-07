# Implementation Plan: Data Backup & Export

## Phase 1: Export
- [x] Task: Build export endpoint and UI
  - [x] Write tests for export JSON structure
  - [x] Collect repos, tags, agents, preferences into single object
  - [x] Add download button to settings panel
  - [x] Include export timestamp and schema version

## Phase 2: Import
- [x] Task: Build import endpoint and UI
  - [x] Write tests for import validation
  - [x] Validate schema version and required fields
  - [x] Support merge (additive) and replace modes
  - [x] Show confirmation with counts before applying

## Phase 3: Polish
- [x] Task: Add error handling and manual verification
  - [x] Write tests for malformed import rejection
  - [x] Clear error messages for validation failures
  - [x] Manual verification: export, clear localStorage, import, verify restored (deferred - manual UI step)
