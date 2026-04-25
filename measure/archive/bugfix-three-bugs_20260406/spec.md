# Specification: Three Critical Bug Fixes

## Track ID
`bugfix-three-bugs_20260406`

## Priority
**P0 — Blocks daily developer workflow.**

## Overview

Three user-reported bugs that severely impact the dashboard's core functionality: models are missing from the agent editor, git pulls fail silently for many repos, and agent logs never display despite log files existing on disk.

---

## Bugs

### BUG-1: Models not populating in Add/Edit Agent dialogs

- **Symptom:** The `/api/models` endpoint returns `{"models":null}`. The model picker dropdown in the agent form shows "No matches" for any search.
- **Root Cause:** `discoverModels()` in `internal/api/agent_handlers.go:336-349` uses `exec.Command("opencode", "models")` without an explicit binary path. The Go process inherits its PATH from however it was launched. If `opencode` is not on the Go process's PATH (it lives at `/home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode`), the command fails silently and returns `nil`.
- **Impact:** Users cannot select a model when creating or editing agents. Without a model, agents are added to crontab with an empty `-m` flag, causing the agent to fail at runtime or use a wrong default.
- **Evidence:** `curl http://localhost:8080/api/models` returns `{"models":null}`. Running `opencode models` from the shell works and returns 252 models.

### BUG-2: Git pulls sporadic/failing with no useful feedback

- **Symptom:** Scheduled pulls fail silently. Manual pulls from the UI show "Failed" with no explanation. 9 out of 24 repos consistently fail to pull.
- **Root Cause (multiple sub-issues):**
  1. **No tracking branch:** 7 repos (advantage-pr, integrated-math-3, ka-math-companion, movie-clips, opencode-configger, primary-advantage, reading-advantage-llm-benchmark) have branches with no upstream tracking set. `git pull --ff-only` fails with "no tracking information."
  2. **Merge conflict:** 1 repo (measure) has unmerged files. `git pull --ff-only` fails with "unresolved conflict."
  3. **Silent failure in scheduler:** `scheduler.RunOnce()` logs the error but doesn't surface it anywhere in the UI.
  4. **No error detail in API:** The `/api/pull` endpoint returns a generic `pull failed: <full error>` string. The frontend shows only "Failed" — the user can't see *why* it failed or which repos are problematic.
- **Impact:** Repos drift out of sync with remotes. The developer doesn't know which repos need attention.
- **Evidence:** Manual test of `git pull --ff-only` on all 24 repos showed 9 failures with specific error messages.

### BUG-3: Agent logs not displaying when clicking Logs button

- **Symptom:** Clicking "Logs" on any agent card shows "No log file found" even though log files exist on disk (confirmed via `ls`).
- **Root Cause:** Agent `log_path` values from the crontab are **relative paths** like `measure/output.log`. The `ReadLogFile` function in `internal/agents/logreader.go` opens the path directly with `os.Open(path)`. Since the Go server's working directory is `/home/daniel-bo/Desktop/dashboard`, it tries to open `/home/daniel-bo/Desktop/dashboard/measure/output.log` instead of `/home/daniel-bo/Desktop/<agent-directory>/measure/output.log`. The file doesn't exist at that path, so `os.Stat` returns "not exist" and the handler returns `{exists: false}`.
- **Impact:** Users cannot view any agent logs. All agents with relative log paths (which is all of them in the current crontab) show "No log file found."
- **Evidence:** `curl http://localhost:8080/api/agents/3/log` returns `{"exists":false,...}` despite `/home/daniel-bo/Desktop/advantage-games/measure/output.log` existing.

---

## Acceptance Criteria

1. `/api/models` returns the full list of available opencode models (200+ entries)
2. The model picker in Add/Edit Agent forms shows models and allows selection
3. Newly created agents have a valid model in their crontab line
4. `/api/pull` returns structured error information including the repo name and error reason
5. The Pull button in the UI shows meaningful error messages (not just "Failed")
6. Repos with tracking issues show a clear explanation of why the pull failed
7. Agent logs display correctly for all agents that have `log_path` set
8. Relative `log_path` values are resolved relative to the agent's `directory`, not the server's working directory
9. All existing tests pass
10. All new code has >80% test coverage

---

## Out of Scope

- Fixing the tracking branch configuration on individual repos (that's a manual git operation)
- Resolving the merge conflict on the measure repo
- Fixing the malformed crontab entries above the automation separator (`cd && && codex`)
- Real-time log streaming
- WebSocket support
