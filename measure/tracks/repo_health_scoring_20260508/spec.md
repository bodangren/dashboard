# Spec: Repository Health Scoring

## Goal
Compute a health score (0-100) for each discovered repository based on staleness, branch hygiene, and remote sync status.

## Acceptance Criteria
- [ ] Health score algorithm defined and tested: considers last commit age, unpushed commits, stale branches, uncommitted changes
- [ ] Score computed during scheduled pull scan
- [ ] Health score displayed on project card with color indicator
- [ ] Sort/filter by health score in the dashboard
- [ ] Health history tracked over time (last 30 days)

## Out of Scope
- Automatic remediation actions
- Health alerts (covered by notification track)
