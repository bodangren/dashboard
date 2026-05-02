# Implementation Plan: Commit Analysis & AI Insights

## Phase 1: Commit Summarizer Service

> Status: Complete

### Tasks

- [x] 1.1: Create `internal/ai/summarizer.go`
  - Interface for LLM provider (abstract for future multi-provider)
  - `SummarizeCommits(repoPath string, commits []CommitInfo) (string, error)`
  - Cache summaries with 5-minute TTL
- [x] 1.2: Add OpenAI-compatible API integration
  - Use existing `OPENAI_API_KEY` env var
  - Implement `openAIProvider` struct
  - Handle rate limiting and retries
- [x] 1.3: Commit event enrichment
  - ActivityEnhancer.EnhanceEvent attaches summary/flags to commit events
  - Wired into gatherCommitEvents via WithActivityEnhancer option

### Verification
- [x] Summarizer can be instantiated without error
- [x] Mock LLM returns predictable output in tests
- [x] Summary appears in /api/activity response for repos with recent commits

---

## Phase 2: Activity Feed Integration

> Status: Complete (backend)

### Tasks

- [x] 2.1: Extend ActivityEvent to include Summary field
  - Add `Summary string` and `Flags []string` to `ActivityEvent` struct
  - Update JSON serialization in api/events.go and ws.ActivityEvent
- [x] 2.2: Background summarization worker
  - ActivityEnhancer uses in-memory cache with TTL
  - Sync summarization on activity request (suitable for low-volume)
- [ ] 2.3: Frontend rendering of summaries
  - Activity feed shows summary badge on summarized events
  - Expandable section for full summary text

### Verification
- [x] Activity feed includes summaries (synchronous enhancement)
- [x] No duplicate summarization requests for same commits (cache)
- [x] Graceful degradation when LLM is unavailable (mock provider fallback)

---

## Phase 3: Issue Detection

> Status: Complete (backend), Pending (frontend)

### Tasks

- [x] 3.1: Detect potential issues in commits
  - Conflict markers in body
  - Rapid successive commits (>5)
  - WIP in message
- [x] 3.2: Issue flag in activity events
  - Add `Flags []string` field to ActivityEvent
  - Populate flags based on detection rules
- [ ] 3.3: Visual indicators in feed
  - Warning icon for flagged events
  - Color-coded badge per flag type

### Verification
- [x] Flagged events have flags populated
- [ ] Flagged events show warning indicator in UI
- [ ] Each flag type is visually distinct in UI

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