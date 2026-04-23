# Track: Fix JS-02 — Deduplicate Service Worker Registration

## Context

Both `app.js` and `diff.js` contain identical service worker registration code (13 lines each). Per tech-debt.md, this should be moved to `utils.js` or a shared module.

## Implementation Plan

### Phase 1 — Extract SW registration to utils.js

- [x] 1.1 Add `registerServiceWorker()` function to `utils.js`
- [x] 1.2 Remove duplicate registration from `app.js`, replace with call to `registerServiceWorker()`
- [x] 1.3 Remove duplicate registration from `diff.js`, replace with call to `registerServiceWorker()`
- [x] 1.4 Verify build succeeds (`go build -o dashboard .`)
- [x] 1.5 Run tests (`go test ./...`)

### Phase 2 — Finalize

- [x] 2.1 Mark JS-02 as Resolved in tech-debt.md
- [x] 2.2 Commit changes