# Specification: Git Dashboard View Enhancements

## Overview

Visual improvements to the existing git commit dashboard: responsive grid layout, commit age badges, and larger project headers.

## Functional Requirements

- **FR-1**: Project cards arranged in a responsive CSS grid — 3 columns on wide viewports (>960px), 2 columns on medium (>600px), 1 column on narrow
- **FR-2**: Each project header row displays a neon green pill/badge showing the relative age of the latest commit (e.g., "2h ago")
- **FR-3**: Project header text size increased to 1.5x the current size

## Non-Functional Requirements

- No changes to the Go backend — all changes are CSS and JS in the static frontend
- Existing functionality (commit rows, hover, click-to-diff) unaffected

## Acceptance Criteria

- Cards display side-by-side at 3 per row on wide screens, 2 on medium, 1 on narrow
- Each project header shows a bold neon green badge with the latest commit's relative age
- Project name and path text are 1.5x larger than current
- Existing functionality (commit rows, hover, click-to-diff) is unaffected

## Out of Scope

- Changes to the diff view
- Backend API changes
