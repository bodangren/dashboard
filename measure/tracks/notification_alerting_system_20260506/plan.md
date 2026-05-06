# Implementation Plan: Notification & Alerting System

## Phase 1: Notification Infrastructure
- [x] Task: Add NotificationService with permission request and localStorage prefs
  - [x] Write tests for permission state management
  - [x] Create service that wraps Notification API with fallback
  - [x] Store preferences (enabled, quiet hours) in localStorage
- [x] Task: Wire health monitoring to notification triggers
  - [x] Write tests for health→notification mapping
  - [x] Trigger notification when health status changes from passing to failing
  - [x] Debounce rapid status changes

## Phase 2: Agent & AI Insight Alerts
- [x] Task: Add agent failure notifications
  - [x] Write tests for agent error detection
  - [x] Parse agent output.log for failure patterns
  - [x] Include last 5 log lines in notification body
- [x] Task: Add AI insight conflict/WIP alerts
  - [x] Write tests for insight flag detection
  - [x] Trigger when AI insights contain "conflict" or "WIP" markers

## Phase 3: UI & Quiet Hours
- [ ] Task: Build notification preferences panel
  - [ ] Write tests for preferences form
  - [ ] Toggle switches for each alert type
  - [ ] Time inputs for quiet hours start/end
- [ ] Task: Manual verification
  - [ ] Verify notifications display in browser
  - [ ] Verify quiet hours suppress alerts
