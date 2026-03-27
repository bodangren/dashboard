# Implementation Plan: MVP — Git Commit Dashboard

## Phase 1 — Project Scaffold & Go Module [checkpoint: 99a3152]

- [x] 1.1 Initialize Go module: `go mod init dashboard`
- [x] 1.2 Create directory structure:
  ```
  cmd/dashboard/main.go
  internal/git/scanner.go
  internal/git/log.go
  internal/git/diff.go
  internal/git/puller.go
  internal/api/handlers.go
  internal/scheduler/scheduler.go
  static/index.html
  static/diff.html
  static/style.css
  static/app.js
  static/diff.js
  ```
- [x] 1.3 Write `main.go` skeleton: wire up HTTP server, embed static files, register routes
- [x] 1.4 **Phase 1 complete**: `go build ./...` succeeds with no errors

## Phase 2 — Git Backend (TDD)

### 2.1 Repo Scanner
- [x] 2.1.1 Write test: `TestScanRepos` — given a temp dir with nested `.git` dirs, scanner returns correct paths
- [x] 2.1.2 Implement `internal/git/scanner.go`: walk directory, collect `.git` parent dirs up to 2 levels deep
- [x] 2.1.3 Run tests: `go test ./internal/git/...`

### 2.2 Git Log Parser
- [x] 2.2.1 Write test: `TestGetCommits` — given a real or temp git repo with commits, returns correct commit structs (hash, message, body, author, timestamp)
- [x] 2.2.2 Implement `internal/git/log.go`: run `git log --format=...` via `os/exec`, parse output into `[]Commit`
- [x] 2.2.3 Run tests

### 2.3 Git Diff Parser
- [x] 2.3.1 Write test: `TestGetDiff` — given a repo and hash, returns unified diff string
- [x] 2.3.2 Implement `internal/git/diff.go`: run `git show <hash>` via `os/exec`, return raw diff
- [x] 2.3.3 Run tests

### 2.4 Git Puller
- [x] 2.4.1 Write test: `TestPullRepo` — verify `git pull` is invoked with correct args (mock exec or use real repo)
- [x] 2.4.2 Implement `internal/git/puller.go`: run `git pull` in repo dir, capture and log errors
- [x] 2.4.3 Run tests

- [x] **Phase 2 complete**: all git package tests pass, `go test ./internal/git/... -cover` ≥ 80%

## Phase 3 — Scheduler

- [ ] 3.1 Write test: `TestScheduler` — verify pull is triggered at correct intervals and runs serially
- [ ] 3.2 Implement `internal/scheduler/scheduler.go`: ticker at 4-hour interval, iterate repos serially with a short sleep between each
- [ ] 3.3 Run tests
- [ ] **Phase 3 complete**: scheduler tests pass

## Phase 4 — API Handlers (TDD)

### 4.1 `/api/projects` endpoint
- [ ] 4.1.1 Write test: `TestProjectsHandler` — mock git layer, assert response JSON shape and sort order
- [ ] 4.1.2 Implement `internal/api/handlers.go`: scan repos, fetch commits, sort by `last_commit_at` desc, return JSON
- [ ] 4.1.3 Run tests

### 4.2 `/api/diff` endpoint
- [ ] 4.2.1 Write test: `TestDiffHandler` — mock git layer, assert diff JSON response
- [ ] 4.2.2 Implement diff handler: validate `repo` and `hash` params, call git diff, return JSON
- [ ] 4.2.3 Run tests

- [ ] **Phase 4 complete**: all API handler tests pass, coverage ≥ 80%

## Phase 5 — Frontend

### 5.1 Dashboard page (`static/index.html` + `static/app.js`)
- [ ] 5.1.1 Write `index.html`: shell with `<div id="projects">`, link `style.css` and `app.js`
- [ ] 5.1.2 Write `app.js`:
  - Fetch `/api/projects` on load
  - Render project cards (repo name, commit rows)
  - Each commit row: message, author, relative timestamp
  - Hover shows commit body (title attribute or inline expand)
  - Click navigates to `diff.html?repo=...&hash=...`
- [ ] 5.1.3 Manual smoke test in browser

### 5.2 Diff page (`static/diff.html` + `static/diff.js`)
- [ ] 5.2.1 Write `diff.html`: shell with metadata block and diff container, "Back" link
- [ ] 5.2.2 Write `diff.js`:
  - Parse `repo` and `hash` from URL params
  - Fetch `/api/diff?repo=...&hash=...`
  - Render commit metadata
  - Render diff lines: `+` prefix → green, `-` prefix → red, else default
- [ ] 5.2.3 Manual smoke test in browser

### 5.3 Styles (`static/style.css`)
- [ ] 5.3.1 Terminal-inspired dark theme: `#111` background, monospace font, green/amber accents
- [ ] 5.3.2 Compact card layout, hover highlight on commit rows
- [ ] 5.3.3 Diff view: green/red line coloring, monospace, scrollable

- [ ] **Phase 5 complete**: dashboard and diff view render correctly in browser

## Phase 6 — Integration & Polish

- [ ] 6.1 Wire scheduler into `main.go`: start on server launch, pass discovered repos
- [ ] 6.2 Embed `static/` into binary with `//go:embed static/*`
- [ ] 6.3 Write integration test: `TestServerEndToEnd` — start server against temp repos, hit `/api/projects`, assert valid response
- [ ] 6.4 Test `go build` → single binary, runs standalone
- [ ] 6.5 Verify on actual `~/Desktop` repos: correct discovery, commits, diff rendering
- [ ] **Phase 6 complete**: binary builds and runs end-to-end

## Phase 7 — Final Review & Commit

- [ ] 7.1 Run full test suite: `go test ./... -cover` — confirm ≥ 80% coverage
- [ ] 7.2 Run `go vet ./...` and `staticcheck ./...` — no issues
- [ ] 7.3 Review against spec acceptance criteria — all checked off
- [ ] 7.4 Commit: `feat(mvp): Implement Git Commit Dashboard`
