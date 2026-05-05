# Implementation Plan: Project Tagging, Groups & Workspaces

## Phase 1: Tag Data Model & Storage
- [ ] Task: Define Tag type and localStorage persistence
  - [ ] Write tests for tag serialization/deserialization
  - [ ] Create Tag type: { repoPath: string, tags: string[] }
  - [ ] Load/save from localStorage with versioning

## Phase 2: Tag UI
- [ ] Task: Add tag chips to project cards
  - [ ] Write tests for tag chip rendering
  - [ ] Render tags as small colored chips below project name
  - [ ] Inline add/edit/delete with Enter to save
- [ ] Task: Add tag filter bar
  - [ ] Write tests for filter behavior
  - [ ] Horizontal bar of active tags at top of repo list
  - [ ] Click to toggle filter; show "Clear" button when active

## Phase 3: Grouping & Polish
- [ ] Task: Add directory grouping option
  - [ ] Write tests for group rendering
  - [ ] Collapsible sections per parent directory
  - [ ] Toggle in preferences panel
- [ ] Task: Manual verification
  - [ ] Verify tags persist across reload
  - [ ] Verify filtering and grouping work together
