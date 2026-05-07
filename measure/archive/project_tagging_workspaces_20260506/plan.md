# Implementation Plan: Project Tagging, Groups & Workspaces

## Phase 1: Tag Data Model & Storage
- [x] Task: Define Tag type and localStorage persistence
  - [x] Write tests for tag serialization/deserialization
  - [x] Create Tag type: { repoPath: string, tags: string[] }
  - [x] Load/save from localStorage with versioning

## Phase 2: Tag UI
- [x] Task: Add tag chips to project cards
  - [x] Write tests for tag chip rendering
  - [x] Render tags as small colored chips below project name
  - [x] Click chip to filter by tag
- [x] Task: Add tag filter bar
  - [x] Write tests for filter behavior
  - [x] Horizontal bar of active tags at top of repo list
  - [x] Click to toggle filter; show "Clear" button when active
- [ ] Task: Inline add/edit/delete tags (deferred)
  - [ ] Add tag via input field on project card
  - [ ] Edit/delete existing tags
  - [ ] Enter to save, Escape to cancel

## Phase 3: Grouping & Polish
- [x] Task: Add directory grouping option
  - [ ] Write tests for group rendering
  - [ ] Collapsible sections per parent directory
  - [ ] Toggle in preferences panel
- [x] Task: Manual verification
  - [ ] Verify tags persist across reload
  - [ ] Verify filtering and grouping work together
