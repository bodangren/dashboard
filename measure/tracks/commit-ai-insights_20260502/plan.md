# Implementation Plan: Commit Analysis & AI Insights

## Phase 1: Commit Summarizer Service

> Status: In Progress

### Tasks

- [x] 1.1: Create `internal/ai/summarizer.go`
  - Interface for LLM provider (abstract for future multi-provider)
  - `SummarizeCommits(repoPath string, commits []git.Commit) (string, error)`
  - Cache summaries with 5-minute TTL
- [x] 1.2: Add OpenAI-compatible API integration
  - Use existing `OPENAI_API_KEY` env var
  - Implement `openaiProvider` struct
  - Handle rate limiting and retries
- [ ] 1.3: Commit event enrichment
  - When new commit events are recorded, queue for summarization
  - On next activity request, attach summary if available

### Verification
- [x] Summarizer can be instantiated without error
- [x] Mock LLM returns predictable output in tests
- [ ] Summary appears in /api/activity response for repos with recent commits

---

## Phase 2: Activity Feed Integration

### Tasks

- [ ] 2.1: Extend ActivityEvent to include Summary field
  - Add `Summary string` to `ActivityEvent` struct in `internal/api/events.go`
  - Update JSON serialization
- [ ] 2.2: Background summarization worker
  - Channel-based job queue for summarization requests
  - Worker pool (1-3 workers) for parallel processing
  - On completion, update event with summary
- [ ] 2.3: Frontend rendering of summaries
  - Activity feed shows summary badge on summarized events
  - Expandable section for full summary text

### Verification
- [ ] Activity feed includes summaries after 30s delay
- [ ] No duplicate summarization requests for same commits
- [ ] Graceful degradation when LLM is unavailable

---

## Phase 3: Issue Detection

### Tasks

- [ ] 3.1: Detect potential issues in commits
  - Large diffs (>500 lines)
  - Rapid successive commits (<5 min apart)
  - Files deleted or renamed
  - Conflict markers in diff
- [ ] 3.2: Issue flag in activity events
  - Add `Flags []string` field to ActivityEvent
  - Populate flags based on detection rules
- [ ] 3.3: Visual indicators in feed
  - Warning icon for flagged events
  - Color-coded badge per flag type

### Verification
- [ ] Flagged events show warning indicator
- [ ] Each flag type is visually distinct
- [ ] No false positives on normal commits

---

## Phase 4: Performance & Polish

### Tasks

- [ ] 4.1: Rate limiting and batching
  - Max 10 LLM calls per minute
  - Batch similar repos if possible
- [ ] 4.2: Summary expiration
  - Keep summaries for 1 hour
  - Re-summarize if commit hash changes
- [ ] 4.3: Error handling polish
  - Log and skip gracefully on LLM failure
  - Show "summary unavailable" in UI

### Verification
- [ ] LLM calls stay within rate limits
- [ ] UI degrades gracefully on API errors
- [ ] All tests pass