# Implementation Plan: Notification & Alerting System

## Phase 1: Notification Infrastructure
- [ ] Task: Add NotificationService with permission request and localStorage prefs
  - [ ] Write tests for permission state management
  - [ ] Create service that wraps Notification API with fallback
  - [ ] Store preferences (enabled, quiet hours) in localStorage
- [ ] Task: Wire health monitoring to notification triggers
  - [ ] Write tests for health→notification mapping
  - [ ] Trigger notification when health status changes from passing to failing
  - [ ] Debounce rapid status changes

## Phase 2: Agent & AI Insight Alerts
- [ ] Task: Add agent failure notifications
  - [ ] Write tests for agent error detection
  - [ ] Parse agent output.log for failure patterns
  - [ ] Include last 5 log lines in notification body
- [ ] Task: Add AI insight conflict/WIP alerts
  - [ ] Write tests for insight flag detection
  - [ ] Trigger when AI insights contain "conflict" or "WIP" markers

## Phase 3: UI & Quiet Hours
- [ ] Task: Build notification preferences panel
  - [ ] Write tests for preferences form
  - [ ] Toggle switches for each alert type
  - [ ] Time inputs for quiet hours start/end
- [ ] Task: Manual verification
  - [ ] Verify notifications display in browser
  - [ ] Verify quiet hours suppress alerts
