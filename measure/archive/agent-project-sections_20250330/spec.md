# Track: Fix Agent Project Sections

## Overview
Fix the agents view to properly organize agents by project sections and ensure read/write functionality works correctly with the new structure.

## Goals
- Revert commits view to original state (removing agent indicators from project cards)
- Modify agents view to show all projects as sections with their associated agents
- Ensure agents are properly grouped under project sections in the crontab
- Auto-set section headers based on project names when creating agents
- Fix log_path field name mapping in frontend

## Acceptance Criteria
- [x] Commits view shows only commits (no agent indicators)
- [x] Agents view shows all projects as sections
- [x] Agents are grouped under their respective project sections
- [x] Creating/editing agents auto-sets section header to project name
- [x] Crontab is reformatted with project sections
- [x] All tests pass
- [x] Build succeeds
