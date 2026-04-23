# Implementation Plan: Git Dashboard View Enhancements

## Phase 1 — Responsive Grid Layout

- [x] 1.1 Update `style.css`: replace single-column `#projects` layout with CSS grid (`grid-template-columns: repeat(auto-fill, minmax(300px, 1fr))`)
- [x] 1.2 Add media queries for 3→2→1 breakpoints (>960px → 3 cols, >600px → 2 cols, ≤600px → 1 col)
- [x] 1.3 Update `main` max-width to fill wider viewport (e.g., `max-width: 1400px`)
- [x] 1.4 Adjust `.project-card` margin for grid gap instead of vertical stacking
- [ ] 1.5 Manual test: verify cards render 3-wide, 2-wide, 1-wide at appropriate breakpoints

## Phase 2 — Project Header Enhancements

- [x] 2.1 Update `style.css`: increase `.project-name` and `.project-path` font sizes to 1.5x (from inherited 13px → ~19.5px for name, ~11.25px → ~17px for path)
- [x] 2.2 Add `.commit-age-badge` CSS class: neon green (`#39ff14`), bold, pill/badge style with padding, border-radius, inline-block in project header
- [x] 2.3 Update `app.js` `renderProject()`: append a badge element to the project header showing `relativeTime(project.commits[0].timestamp)`
- [x] 2.4 Manual test: verify badge appears in header, text is larger, grid still works

## Phase 3 — Polish & Verify

- [x] 3.1 Verify commit rows, hover tooltips, and click-to-diff still work correctly in grid layout
- [x] 3.2 Verify diff view is unaffected
- [x] 3.3 Test with many repos (scroll behavior, no overflow issues)
- [x] 3.4 Commit: `feat(ui): responsive grid layout, commit age badges, larger headers`
