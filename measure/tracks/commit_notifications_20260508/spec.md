# Spec: Commit and Agent Failure Notifications

## Goal
Add a notification system to the dashboard that alerts the user about important events: new commits in watched repos, agent run failures, and long-running agents.

## Acceptance Criteria
- [ ] Backend notification rules engine evaluates events against user-defined rules
- [ ] WebSocket push delivers notifications to connected clients in real time
- [ ] Frontend notification bell with unread count and history panel
- [ ] Rules: new commit in repo X, agent failure rate > threshold, agent run duration > threshold
- [ ] Notifications stored in SQLite with read/unread state
- [ ] Mark-all-read and dismiss individual notifications

## Out of Scope
- Email or push notifications (browser only)
- Mobile notification support
