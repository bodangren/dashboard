# Track: Data Backup & Export

## Overview
Allow exporting dashboard data, settings, and agent configurations for backup or migration.

## Goals
- Export all dashboard data as JSON
- Import previously exported data
- Include repos, tags, agent configs, and preferences

## Acceptance Criteria
- [ ] Export button generates a downloadable JSON file
- [ ] JSON contains repos, tags, agents, and preferences
- [ ] Import validates schema and merges or replaces data
- [ ] Tests cover export format and import validation

## Non-Goals
- Cloud sync
- Automatic scheduled backups
