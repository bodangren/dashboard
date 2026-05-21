# Track: Agent Performance Metrics & Analytics

## Overview
Expose quantitative performance data for every configured agent, enabling users to identify slow, flaky, or expensive automation at a glance. Metrics are computed from execution history stored in SQLite and surfaced through a new dashboard analytics view.

## Goals
- Track per-agent run duration, success rate, and last-run timestamp
- Track per-project agent cost (total runs, average duration)
- Surface trends over 7-day and 30-day windows
- Highlight agents with >20% failure rate or avg duration >5min

## Acceptance Criteria
- [ ] `GET /api/agents/metrics` returns aggregated stats per agent
- [ ] `GET /api/agents/metrics?project=<path>` returns stats scoped to a project
- [ ] Metrics include: total_runs, success_rate_pct, avg_duration_ms, p95_duration_ms, last_run_at
- [ ] New "Metrics" tab in agent panel with sortable table
- [ ] Agents exceeding thresholds render with warning badges
- [ ] All aggregation logic covered by unit tests
- [ ] Zero impact on agent execution hot path (async aggregation only)

## Non-Goals
- Real-time metric streaming (batch aggregation is fine)
- Cost estimation in monetary terms
- Historical graph visualizations (sparklines can be added later)
- Alerting on threshold breaches (use existing notification system)
