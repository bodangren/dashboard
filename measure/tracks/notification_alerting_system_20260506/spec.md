# Track: Notification & Alerting System

## Overview
Turn the passive dashboard into an active assistant by adding browser/desktop notifications for critical events.

## Goals
- Notify when repo health goes critical
- Notify when agent execution fails
- Notify when AI insights detect merge conflicts or WIP markers
- Configurable thresholds and quiet hours

## Acceptance Criteria
- [ ] Browser notifications request permission and display
- [ ] Health alerts fire when checks fail
- [ ] Agent failure alerts include project name and last log lines
- [ ] Quiet hours configurable in UI
- [ ] All notification logic covered by unit tests

## Non-Goals
- Email or Slack integrations
- Mobile push notifications
