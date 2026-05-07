# Track: Project Tagging, Groups & Workspaces

## Overview
Allow custom tags, labels, and auto-grouping to organize the flat repo list.

## Goals
- Custom tags/labels (e.g., "client", "archive", "urgent")
- Filter repos by tag
- Auto-group by parent directory
- Persist preferences in localStorage

## Acceptance Criteria
- [ ] Tags can be added, edited, and removed per repo
- [ ] Main view can filter by active tag
- [ ] Optional auto-group by parent directory
- [ ] Preferences persist across sessions
- [ ] Tests cover tag CRUD and filtering logic

## Non-Goals
- Server-side persistence
- Shared tags across users
