# Track: Advanced Search with Server-Side Indexing

## Overview
Implement server-side indexing and search across all repositories to find code snippets or commit messages without opening individual projects.

## Goals
- Index commit messages, file names, and diff content
- Fast fuzzy search across all repos
- Search results show repo, commit, and matched snippet

## Acceptance Criteria
- [ ] Background indexer scans all repos on startup and after pulls
- [ ] Fuzzy search endpoint returns ranked results
- [ ] Results page shows repo, commit hash, message, and snippet
- [ ] Index updates incrementally after new commits
- [ ] Tests cover indexing and search logic

## Non-Goals
- Full-text code content indexing (file names and commits only)
- Elasticsearch or external search engine
