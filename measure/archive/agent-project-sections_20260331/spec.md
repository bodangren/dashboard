# Track: Agent Project Sections — Crontab Reorganization

## Overview

The crontab's automation section must be structurally organized: a separator line marks where automation begins, every discovered project gets a section header (even with zero agents), and all agents for a project sit under its header. Currently `AddAgent` appends at the end (creating duplicate headers), `DeleteAgent` removes headers prematurely, and `UpdateAgent` ignores section headers entirely.

## Desired Crontab Structure

```
# <preamble — env vars, maintenance cron jobs, untouched>

# ==AUTOMATION BELOW THIS LINE==
# ~/Desktop/dashboard
0 */4 * * * cd /home/daniel-bo/Desktop/dashboard && opencode -m ... run ...

# ~/Desktop/bus-math-v2
45 0,6,12,18 * * * cd /home/daniel-bo/Desktop/bus-math-v2 && opencode -m ... run ...

# ~/Desktop/verbal

# ~/Desktop/kanban-measure
15 1,5,9,13,17 * * * cd /home/daniel-bo/Desktop/kanban-measure && codex -m ... run ...
```

Rules:
1. Everything above the separator is **read-only** — never modified by agent CRUD.
2. The separator line is `# ==AUTOMATION BELOW THIS LINE==`. Created on first use if absent.
3. Below the separator: one `# <project-dir>` header per project, followed by its agents (zero or more).
4. Section headers use the project's full path (e.g. `# ~/Desktop/dashboard`) for unambiguous identification.
5. Empty project sections (no agents) still appear — they're placeholders for future agents.

## Goals

- `ReorganizeAutomation(projects []string)` rebuilds everything below the separator: creates headers for all projects, preserves existing agents under correct headers, removes stale empty sections.
- `AddAgent` inserts an agent under its project's section header (not at end of file).
- `DeleteAgent` removes the agent; removes the section header only if the project has no remaining agents AND the project is not in the known projects list.
- `UpdateAgent` updates the agent line and its section header comment if the project changed.
- The `AgentHandler` receives the project list (same repos as `Handler`) so it can call reorganize after mutations.
- Frontend passes the project list on create/edit so the backend can maintain correct sections.
- Malformed agent lines (empty directory from old crontab bugs) are handled gracefully.

## Acceptance Criteria

- [ ] Separator line `# ==AUTOMATION BELOW THIS LINE==` exists in crontab after first agent operation
- [ ] Every project from `~/Desktop` has a `# <path>` section header below the separator
- [ ] Adding an agent places it under the correct project header, never creates duplicate headers
- [ ] Deleting the last agent for a project removes the section header
- [ ] Deleting a non-last agent preserves the section header
- [ ] Editing an agent to a different project moves it to the new project's section
- [ ] Toggling an agent does not affect section structure
- [ ] Content above the separator is never modified
- [ ] `Crontab.String()` output is deterministic and parseable round-trip
- [ ] All existing tests pass
- [ ] New tests cover: reorganize, add-under-section, delete-keeps-header, move-across-projects
