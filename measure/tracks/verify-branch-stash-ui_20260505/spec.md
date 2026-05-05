# Track: Verify Branch & Stash UI

## Context

Track: verify-branch-stash-ui_20260505

Per autonomous_prompt.md Section 2.1, the previous phase (platform-git-features_20260505) was marked complete but has a known open item: "FRONTEND-03: Branch/stash UI needs user testing". This track implements the mandatory functionality review before declaring the previous phase truly done.

## Spec

Verify that branch and stash UI actually works in the browser:
1. Branch badge shows on project cards
2. Branch dropdown opens and shows branch list
3. Branch checkout works (POST /api/branches/checkout)
4. Stash toggle button opens stash panel
5. Stash list renders correctly
6. Apply stash works
7. No console errors during interaction
8. Visual layout renders correctly

## Tasks

### Phase 1: Browser Verification

- [ ] Start dev server
- [ ] Connect via browser-harness
- [ ] Verify branch badge visible on project card
- [ ] Verify branch dropdown opens on click
- [ ] Verify stash panel opens on toggle
- [ ] Test branch checkout interaction
- [ ] Test stash apply interaction
- [ ] Check for console errors
- [ ] Take screenshot for verification
- [ ] Fix any issues found
- [ ] Stop dev server
- [ ] Commit