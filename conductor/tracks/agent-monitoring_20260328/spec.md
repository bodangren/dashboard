# Specification: Autonomous Agent Monitoring

## Overview

Add a new "Agents" tab/page to the dashboard that reads the user's crontab, parses autonomous agent entries (opencode, gemini, codex), and presents them in a manageable view with full CRUD + toggle capabilities.

## Functional Requirements

### FR-1: Agent Parsing

- Parse `crontab -l` output to extract agent entries, distinguishing agent lines from comments, env vars, and other cron jobs
- Extract per-agent: schedule (cron expression), working directory, harness type (opencode/gemini/codex), model, prompt/command, log path
- Detect commented-out entries as "disabled" agents (preserving the schedule and config)
- Identify agent lines by the presence of known harness commands: `opencode`, `gemini`, `codex`

### FR-2: Agent List View

- Separate page (`agents.html`) accessible via navigation tab from the git dashboard
- Group agents by project directory
- Each agent card shows: project name, harness label, model, schedule (human-readable), enabled/disabled status, log path
- Terminal-styled dark theme consistent with existing dashboard

### FR-3: Agent CRUD

- **Create**: Form to add a new agent — select directory (from discovered git repos), harness (opencode/gemini/codex), model name, cron schedule, prompt file/command, log path
- **Edit**: Inline or modal edit of any existing agent field
- **Delete**: Remove an agent entry from crontab (with confirmation prompt)
- **Toggle**: Enable/disable an agent by commenting/uncommenting its crontab line

### FR-4: Schedule Editor

- Human-friendly schedule input that generates cron expressions
- Display raw cron expression alongside human-readable description
- Support common patterns: "Every N hours", "Daily at HH:MM", "Times: HH,HH,HH"

### FR-5: Log Viewer

- Expandable section per agent showing last N lines of its log file
- Indicate last run time based on log file modification time
- Handle missing log files gracefully

## Non-Functional Requirements

- Preserving non-agent crontab lines (env vars, update jobs, comments) untouched when writing back
- Crontab writes are atomic (read full, modify agent lines only, write full back)
- All changes are frontend + Go backend API additions — no external dependencies

## Acceptance Criteria

- Dashboard nav includes "Agents" tab linking to agent monitoring page
- All existing agent entries from crontab are parsed and displayed correctly
- Commented-out agents shown as disabled with toggle to re-enable
- Creating a new agent adds a valid crontab entry
- Editing an agent updates its crontab line in-place
- Deleting an agent removes its crontab line (with confirmation)
- Non-agent crontab lines (env vars, update jobs) are preserved through all operations
- Log file contents viewable per agent

## Out of Scope

- Real-time agent process monitoring (only crontab + log inspection)
- Agent output parsing or structured log analysis
- Authentication or multi-user support
- Claude Code harness support (deferred to future track)
