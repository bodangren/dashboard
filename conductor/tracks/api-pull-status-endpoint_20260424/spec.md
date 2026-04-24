# Specification: `/api/pull/status` GET Endpoint

## Overview

Implement the GET endpoint for `/api/pull/status` to return the status of background pull operations for all repositories.

## Functional Requirements

- [ ] `GET /api/pull/status` returns JSON with per-repository pull status
- [ ] Response includes: repo path, last pull timestamp, any errors from last pull
- [ ] Track in-progress pull operations
- [ ] Return appropriate HTTP status codes (200 OK, 500 on error)

## Data Model

```json
{
  "statuses": [
    {
      "repo": "/path/to/repo",
      "lastPullTime": "2026-04-24T12:00:00Z",
      "lastError": "",
      "inProgress": false
    }
  ]
}
```

## Acceptance Criteria

1. Endpoint returns valid JSON
2. All configured repositories are listed
3. In-progress pulls are tracked and reflected in response
4. Last error is surfaced when present
5. Works alongside existing POST `/api/pull` endpoint

## Out of Scope

- Push status (read-only for now)
- Per-branch status