# Specification: Improved Search & Filtering

## Overview

Implement server-side search and filtering capabilities for the dashboard, allowing users to find specific code snippets, commit messages, or project content without opening individual repositories.

## Functional Requirements

1. **Search Index Building**
   - Build an in-memory search index of all repository content on startup
   - Index commit messages, file names, and optionally file contents
   - Update index incrementally when repos pull new commits

2. **Search API Endpoint**
   - `GET /api/search?q=<query>` endpoint for full-text search
   - Support filtering by repository, date range, author
   - Return ranked results with context snippets

3. **Filter UI**
   - Search input in the main dashboard header
   - Filter options: repo selector, date range picker, author filter
   - Real-time search-as-you-type with debouncing

4. **Results Display**
   - Display matching commits in a collapsible panel
   - Show match context with highlighted terms
   - Click to navigate to full commit/diff view

## Non-Functional Requirements

- Search should respond in <200ms for typical queries
- Index should not significantly increase memory footprint
- Incremental updates should not block main goroutine

## Acceptance Criteria

- [ ] Search endpoint returns relevant results within 200ms
- [ ] Search UI shows results as user types (debounced)
- [ ] Filters correctly narrow results
- [ ] Clicking result navigates to commit detail view
- [ ] Index builds on startup without blocking server start

## Out of Scope

- Full-text content search (only commit messages and file names)
- Search result caching
- Advanced query syntax (AND/OR/NOT operators)
- Search history