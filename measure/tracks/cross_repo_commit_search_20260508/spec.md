# Track: Cross-Repository Commit Search & Filtering

## Overview
Add powerful search and filtering capabilities to the unified commit feed, allowing users to find commits across all discovered repositories by message content, author, or date range. This complements the existing server-side search index with a focused commit-centric query interface.

## Goals
- Search commits across all repos by message substring or regex
- Filter commits by author name or email
- Filter commits by date range (relative: "last 7 days", or absolute)
- Combine filters (author + date + message)
- Persist recent searches in localStorage

## Acceptance Criteria
- [ ] `GET /api/commits/search` endpoint accepts `q`, `author`, `since`, `until`, and `repo` query params
- [ ] Endpoint returns paginated commit results ordered by timestamp desc
- [ ] Search is case-insensitive and matches commit message subject and body
- [ ] Date filters accept ISO dates and relative shortcuts (`1d`, `7d`, `30d`)
- [ ] Frontend search bar supports all filter modes with clear UX
- [ ] All handler and service logic covered by unit tests
- [ ] Performance: <100ms for scans up to 10k commits

## Non-Goals
- Full-text search inside diff contents (use existing `/api/search`)
- Elasticsearch or external indexing (SQLite queries are sufficient)
- Saved searches or search sharing
