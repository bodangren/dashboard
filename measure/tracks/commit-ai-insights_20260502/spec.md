# Specification: Commit Analysis & AI Insights

## Overview

Integrate a local LLM to provide summaries of recent changes across all repositories, identifying potential bugs or architectural regressions directly in the dashboard.

## Goals

- Display AI-generated summaries of recent commits in the activity feed
- Highlight potential issues (large diffs, rapid changes, conflict markers)
- Keep LLM calls efficient via caching and batching

## Technical Approach

- Use the existing ActivityFeed WebSocket infrastructure
- Add a new `commit_analyzer` package with LLM integration
- Cache commit summaries to avoid redundant LLM calls
- Display summaries inline in the activity feed alongside regular events

## Out of Scope

- Self-hosted LLM deployment (use external API)
- Detailed code review (just summaries/highlights)
- Real-time commit scanning (polling-based, every 5 min)