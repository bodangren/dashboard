# Specification: Enhanced Agent Orchestration & Monitoring

## Overview

Expand the `/api/agents` capability to include real-time log streaming via WebSockets, manual trigger overrides for cron jobs, and better error reporting for failed agent runs.

## Goals

1. **Real-Time Log Streaming**: WebSocket-based live updates for agent activity
2. **Manual Trigger Override**: Allow users to force-run any cron job on demand
3. **Error Reporting**: Capture and display stderr/exit codes for failed runs

## Technical Approach

### Real-Time Events
- Reuse existing `logHub` WebSocket hub
- Add `/ws/activity` endpoint for activity feed subscriptions
- Broadcast events as JSON: `{id, type, repo, message, timestamp, metadata}`

### Manual Trigger
- `POST /api/agents/<id>/trigger` endpoint
- Returns `202 Accepted` immediately
- Queues agent for async execution
- WebSocket notification on completion

### Error Capture
- AgentStateMap stores `{exitCode, stderr, lastRun}` per agent
- Error badge on UI when exitCode != 0
- Modal to view full error with log context

## API Changes

```
GET  /api/agents          - list all agents (existing)
POST /api/agents          - create agent (existing)
GET  /api/agents/<id>     - get agent details (existing)
PUT  /api/agents/<id>     - update agent (existing)
DELETE /api/agents/<id>    - delete agent (existing)
POST /api/agents/<id>/trigger - manual trigger (NEW)
GET  /ws/activity          - WebSocket subscription (NEW)
```

## File Changes

- `internal/api/agent_handlers.go` - add trigger endpoint
- `internal/ws/hub.go` - existing hub supports activity broadcast
- `static/agents.html` - add trigger button
- `static/agents.js` - handle trigger + WebSocket subscription