# Repository Health Monitoring

## Overview

Add a repository health monitoring feature that provides at-a-glance status indicators for each repository beyond just commit history. This addresses a product vision gap where users currently see only recent commits but lack visibility into repo health signals that would help prioritize maintenance work.

## Problem Statement

Solo developers managing multiple repositories on ~/Desktop have no easy way to identify repos that need attention—those with uncommitted changes, stale branches, or diverged remotes. Currently, developers must manually check each repo to assess its health status.

## Solution

Implement a health monitoring system that scans repositories for key health indicators and displays them prominently in the project cards:

### Health Indicators
- **Dirty Working Tree**: Show uncommitted changes count (modified, staged, untracked files)
- **Branch Divergence**: Detect when local branch is ahead/behind remote
- **Stale Branches**: Count branches older than 30 days without activity
- **Empty Repos**: Repos with zero commits (already partially handled)

### Display
- Add a health status badge to project cards (green/yellow/red)
- Expandable health details section on hover or click
- Color-coded indicators: green (healthy), yellow (minor issues), red (needs attention)

## Acceptance Criteria

- [ ] Health status endpoint returns per-repo health data
- [ ] Project cards display health indicator badges
- [ ] Health details expandable on interaction
- [ ] Performance: health checks complete within 5 seconds for 50 repos
- [ ] Tests cover health indicator calculations and edge cases

## Out of Scope

- Historical health tracking over time
- Automated cleanup actions (delete stale branches, etc.)
- External CI/CD integration status
