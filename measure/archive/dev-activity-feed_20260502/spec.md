# Track: Unified Developer Activity Feed

## Overview

Create a centralized timeline view that combines Git commits, Agent logs, and CI/CD status into a single unified stream. This gives developers a single pane to see all development activity across their repositories.

## Functional Requirements

- **Activity Feed UI**: New `/activity` page with a chronological timeline of all development events
- **Event Types**: Commits, Agent runs (with status/failure indicators), Pull status updates
- **Filtering**: Filter by event type, repository, time range
- **Real-time Updates**: New events appear via WebSocket without page refresh
- **Responsive Design**: Works on mobile and desktop

## Non-Functional Requirements

- Timeline loads within 2s for 100+ events
- Efficient aggregation: Query each source in parallel
- Persist last-seen timestamp to avoid duplicate events on reconnect

## Acceptance Criteria

1. `/activity` route renders a unified timeline
2. Each event shows: type icon, repo name, description, timestamp
3. Filter controls allow toggling event types on/off
4. New events stream in via existing WebSocket infrastructure
5. Activity persists across page refreshes (last-read cursor)

## Out of Scope

- Full-text search within events
- Export to CSV/JSON
- User mentions or notifications
- Integration with external CI systems (GitHub Actions, etc.)