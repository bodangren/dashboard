# Plan: Repository Health Scoring

## Phase 1: Score Algorithm (TDD)
- [ ] Write tests for health score calculation from repo metadata
- [ ] Implement `CalculateHealthScore` pure function
- [ ] Validate edge cases (empty repo, very old repo, clean repo)

## Phase 2: Backend Integration (TDD)
- [ ] Write tests for health score persistence in repo record
- [ ] Add health_score and health_history fields to repo storage
- [ ] Wire score computation into scheduled pull scan

## Phase 3: Frontend Display (TDD)
- [ ] Write tests for HealthBadge component
- [ ] Implement HealthBadge with color coding
- [ ] Write tests for health sort/filter in project grid
- [ ] Add sort by health score and filter by health range

## Phase 4: History Chart (TDD)
- [ ] Write tests for HealthHistory sparkline rendering
- [ ] Implement 30-day health history mini-chart on project card

## Phase 5: Integration & Verification
- [ ] Verify all repos get scored during next scan
- [ ] Update tracks.md and commit
