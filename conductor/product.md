# Product Definition

## Overview

**Git Commit Dashboard** — A personal local web app that auto-scans `~/Desktop` for git repositories and displays their latest commits in a unified view. Designed for solo developer use.

## Target Users

- Solo developer (the machine owner)

## Core Features

- **Auto-discovery**: Scan `~/Desktop` for all git repositories on startup
- **Scheduled pulls**: Run `git pull` on each repo 6 times per day, serially with wait time between each pull, to stay up to date with remote
- **Commit feed**: Display latest commits per repo — commit message, author, and timestamp
- **Dynamic sorting**: Project cards sorted by most recently committed repo (floats to top)
- **Commit hover**: Mouse over a commit to view the associated commit note/body
- **Diff report**: Click a commit to open a detailed diff report showing files changed, insertions, and deletions
- **Local web access**: Served from localhost, accessed in any browser

## Non-Goals

- Multi-user or networked access
- Manual project path configuration
- Push or write operations to repositories
