# Track: Multi-Platform Support & Advanced Git Features

## Context

The roadmap item "Multi-Platform Support & Advanced Git Features" was marked as TODO in tracks.md. The product definition specifies a personal local web app for solo developer use. This track implements branch management and stash viewing features that enhance the existing git operations without introducing multi-user authentication (which contradicts the solo-developer product definition).

## Implementation Plan

### Phase 1: Branch Management

- [ ] Add branch list API endpoint (GET /api/branches?repo=<path>)
- [ ] Add branch checkout handler (POST /api/branches/checkout with JSON body)
- [ ] Add create branch handler (POST /api/branches with JSON body)
- [ ] Add delete branch handler (DELETE /api/branches)
- [ ] Write tests for branch handlers
- [ ] Add frontend branch UI in project cards (show current branch badge, dropdown for switch/create/delete)

### Phase 2: Stash Management

- [ ] Add stash list API endpoint (GET /api/stash?repo=<path>)
- [ ] Add stash apply handler (POST /api/stash/apply)
- [ ] Add stash drop handler (DELETE /api/stash/<stash-id>)
- [ ] Write tests for stash handlers
- [ ] Add stash viewer UI accessible from project card

### Phase 3: UI Integration & Polish

- [ ] Integrate branch picker into project card header
- [ ] Add stash panel to diff/activity view
- [ ] Verify all tests pass
- [ ] Manual smoke test in browser