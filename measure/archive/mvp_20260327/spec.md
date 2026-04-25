# Specification: MVP — Git Commit Dashboard

## Goal

Build a local web app (Go backend + vanilla HTML/CSS/JS frontend) that auto-discovers git repositories under `~/Desktop`, keeps them up to date via scheduled pulls, and presents their latest commits in a terminal-styled dashboard.

## Acceptance Criteria

### Repo Discovery
- [ ] On startup, the server walks `~/Desktop` and finds all directories containing a `.git` folder (up to 2 levels deep)
- [ ] Discovered repos are stored in memory; no config file required

### Scheduled Pulls
- [ ] `git pull` runs on each discovered repo **6 times per day** (every 4 hours)
- [ ] Pulls run **serially** with a short wait between each repo (avoid hammering remotes)
- [ ] Pull errors are logged server-side but do not crash the server

### Commit Feed API
- [ ] `GET /api/projects` returns JSON array of projects, each with:
  - `name` (repo directory name)
  - `path` (absolute path)
  - `last_commit_at` (timestamp of most recent commit)
  - `commits` — array of last N (default 10) commits, each with:
    - `hash` (short 7-char hash)
    - `message` (subject line)
    - `body` (full commit body/notes, may be empty)
    - `author`
    - `timestamp` (ISO 8601)
- [ ] Projects sorted by `last_commit_at` descending (most recent first)

### Diff API
- [ ] `GET /api/diff?repo=<path>&hash=<hash>` returns JSON with:
  - `hash`, `message`, `author`, `timestamp`
  - `diff` — raw unified diff string for the commit

### Frontend — Dashboard
- [ ] Page loads and fetches `/api/projects`, renders one card per repo
- [ ] Cards are sorted: most recently committed repo at top
- [ ] Each card shows: repo name, list of recent commits (message, author, relative timestamp)
- [ ] Hovering a commit row shows the commit body/note (tooltip or inline expand)
- [ ] Clicking a commit row navigates to the diff view for that commit

### Frontend — Diff View
- [ ] Displays commit metadata (hash, message, author, timestamp)
- [ ] Renders unified diff with `+` lines in green, `-` lines in red
- [ ] "Back" button/link returns to the dashboard

### Visual Style
- [ ] Terminal-inspired dark theme (`#111` background, monospace font, green/amber accents)
- [ ] Compact, information-dense layout
- [ ] No external CSS frameworks or JS libraries

### Server
- [ ] Single `go build` produces a self-contained binary
- [ ] Serves on `http://localhost:8080` by default
- [ ] Static frontend files embedded in the binary (using `go:embed`)

## Out of Scope

- Authentication or multi-user support
- Push / write operations to repos
- Manual repo path configuration
- External databases
- Build step for frontend (no bundler)
