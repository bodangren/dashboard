# Agent Editor Fix — Specification

## Problem

The Agents tab cannot correctly edit agent timing/configuration because:

1. **Wrong CLI format**: Parser expects `--model` and `--prompt` flags, but real opencode entries use `-m` for model and `run <path>` for prompts
2. **Wrong writer output**: `buildAgentLine` generates `--model`/`--prompt` instead of `-m`/`run`
3. **No section headers**: Comment lines above agent lines are not captured as section headers for grouping
4. **Multi-harness UI**: UI shows opencode/gemini/codex but user wants opencode-only
5. **Wrong log redirect regex**: Parser expects `>>` but real entries use `>`
6. **Full binary path not preserved**: Real crontab uses full path to opencode binary

## Scope

- Fix parser regexes to match real opencode `run` subcommand format: `opencode -m <model> run <prompt>`
- Support both `>` and `>>` redirect styles
- Capture full binary path for opencode
- Add section header concept: comment line immediately above an agent line becomes its `SectionHeader`
- Fix writer to emit correct opencode format
- Restrict harness to `opencode` only in UI
- Add model dropdown populated from `opencode models` (or a static model list fetched at load)
- Fields in form: Schedule, Project Path, Model (dropdown), Prompt Path, Log Output Path

## Acceptance Criteria

- [ ] Parser correctly extracts model from `-m provider/model` flag
- [ ] Parser correctly extracts prompt path from `run <path>` positional
- [ ] Parser handles `>` redirect in addition to `>>`
- [ ] Section header comments (line above agent line) are captured in the Agent struct
- [ ] Writer generates lines like: `30 3,7,11,15,23 * * * cd /path/to/project && /full/path/opencode -m provider/model run prompt/path > /log/path 2>&1`
- [ ] Adding an agent inserts a comment section header line above it
- [ ] UI only shows opencode harness
- [ ] Model field is a searchable dropdown of available models
- [ ] All existing tests pass + new tests for section headers, correct CLI format
- [ ] `CI=true bun run test --run` passes
- [ ] `cd app && npm run build` passes (if applicable)
