# Implementation Plan: JS-01 — Replace `var` with `const`/`let`

## Phase 1 — Convert var to const/let in agents.js

### 1.1 Analyze and convert var declarations

- [x] 1.1.1 Count all `var` occurrences in agents.js (baseline) — 93 found
- [x] 1.1.2 For each `var`, determine if the variable is reassigned — 9 identified (loop counters: i×3, h, d; cachedModels, cachedRepos, hourCheckboxes, dayCheckboxes)
- [x] 1.1.3 Replace `var` with `const` for non-reassigned variables — 93 total, 84 became const
- [x] 1.1.4 Replace `var` with `let` for reassigned variables — 9 changed to let
- [x] 1.1.5 Verify all `var` replaced — count is 0

### 1.2 Verify

- [ ] 1.2.1 Open agents.html in browser, verify page loads without console errors (manual)
- [x] 1.2.2 Conductor - Mark JS-01 as Resolved in tech-debt.md

## Phase 2 — Final commit

- [ ] 2.1.1 Commit: `fix(js): replace var with const/let in agents.js per JS-01`
