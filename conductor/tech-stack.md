# Tech Stack

## Backend

- **Language**: Go
- **Role**: HTTP server, git repo discovery, scheduled `git pull`, git log/diff parsing
- **Key packages**:
  - `net/http` — serve the dashboard and API endpoints
  - `os/exec` — run git commands (`git pull`, `git log`, `git diff`)
  - `path/filepath` — walk `~/Desktop` to discover repos
  - `encoding/json` — API responses
  - `time` — scheduling pulls (6x/day, serially)

## Frontend

- **Stack**: Vanilla HTML, CSS, JavaScript (no framework, no build step)
- **Served**: As static files embedded in or served by the Go binary
- **JS role**: Fetch JSON from Go API endpoints, render project cards, handle hover/click interactions

## Data Flow

1. Go server scans `~/Desktop` for `.git` dirs on startup
2. Scheduler runs `git pull` on each repo sequentially, 6 times/day
3. On page load, frontend fetches `/api/projects` → list of repos with latest commits
4. On commit click, frontend fetches `/api/diff?repo=...&hash=...` → diff data
5. All data rendered client-side from JSON

## Development

- **Build**: `go build` → single binary
- **Run**: `./dashboard` → serves on `http://localhost:8080` (or configurable port)
- **No database**: all data read live from git on request or cached in memory

## Constraints

- Read-only access to repos (no push, no write)
- No external services or network dependencies beyond `git pull`
