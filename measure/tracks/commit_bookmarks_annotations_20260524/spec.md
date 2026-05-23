# Spec: Commit Bookmarking & Personal Annotations

## Problem

The dashboard displays a high-volume commit feed across many repositories, but there is no way for the solo developer to mark commits as important, add quick context, or revisit them later. Valuable commits (e.g., critical fixes, architectural decisions, AI-flagged insights) scroll out of view and are lost in the noise. The developer must rely on memory or external note-taking tools.

## Goals

1. Allow users to bookmark any commit from the commit feed or diff view.
2. Attach short personal annotations (notes) to bookmarked commits.
3. Provide a dedicated "Bookmarks" view to browse, search, and filter bookmarked commits.
4. Surface bookmark status and annotations in the diff view and activity feed.
5. Persist bookmarks and annotations in SQLite with the existing local-server architecture.

## Non-Goals

- Public or shared bookmarks (solo use only).
- Rich text or markdown in annotations (plain text only).
- Git notes integration (store separately from git).
- Email or notification reminders for bookmarks.

## Approach

- Add a `bookmarks` table: `id`, `repo`, `hash`, `note`, `created_at`, `updated_at`.
- Add a `BookmarkStore` Go package with CRUD operations.
- Extend the frontend commit cards and diff view with a bookmark toggle (star icon) and an inline note editor.
- Add a `/bookmarks` page listing all bookmarks with repo filter and search.
- Wire bookmark counts into the existing project cards (e.g., "3 bookmarks" badge).
- Reuse existing API patterns (HandlerConfig, JSON responses) and frontend patterns (vanilla JS, CSS custom properties).

## Success Criteria

- User can bookmark/unbookmark a commit from the commit feed and diff view.
- User can add, edit, and delete a personal note on any bookmarked commit.
- The `/bookmarks` page lists all bookmarks with full-text search and repo filter.
- Bookmarked commits are visually distinguished in the commit feed (star icon).
- All data persists across server restarts (SQLite).
- New code has Go unit tests for the store and handler; frontend is manually verified.
