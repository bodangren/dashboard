# Track: Commit Activity Calendar and Streak Tracking

## Overview

Add a GitHub-style contribution calendar heatmap to the dashboard that visualizes commit activity across all discovered repositories. Include streak tracking (current and longest) to gamify developer productivity.

## Goals

1. Generate commit activity data for the last 365 days across all repos
2. Render a color-coded calendar heatmap (7-day columns, 52 weeks)
3. Display current streak and longest streak metrics
4. Allow clicking a day to see commits from that date
5. Update incrementally as new commits arrive

## Non-Goals

- GitHub API integration (local repos only)
- Comparative analytics against other users
- Predictive forecasting

## Success Criteria

- Calendar renders correctly on desktop and mobile
- Streak calculation is accurate (consecutive days with ≥1 commit)
- Clicking a day navigates to filtered commit view
- Performance: <100ms to generate data for 20 repos
- No external dependencies beyond existing frontend stack
