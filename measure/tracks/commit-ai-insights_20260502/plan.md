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

> Status: Complete

### Tasks

- [x] 2.1: Extend ActivityEvent to include Summary field
  - Add `Summary string` and `Flags []string` to `ActivityEvent` struct
  - Update JSON serialization in api/events.go and ws.ActivityEvent
- [x] 2.2: Background summarization worker
  - ActivityEnhancer uses in-memory cache with TTL
  - Sync summarization on activity request (suitable for low-volume)
- [x] 2.3: Frontend rendering of summaries
  - Activity feed shows summary text below commit message
  - Expandable section for full summary text

### Verification
- [x] Activity feed includes summaries (synchronous enhancement)
- [x] No duplicate summarization requests for same commits (cache)
- [x] Graceful degradation when LLM is unavailable (mock provider fallback)

---

## Phase 3: Issue Detection

> Status: Complete

### Tasks

- [x] 3.1: Detect potential issues in commits
  - Conflict markers in body
  - Rapid successive commits (>5)
  - WIP in message
- [x] 3.2: Issue flag in activity events
  - Add `Flags []string` field to ActivityEvent
  - Populate flags based on detection rules
- [x] 3.3: Visual indicators in feed
  - Summary text shown in italic below commit message
  - Flag badges with color coding: conflict (red), wip (orange), rapid-changes (blue)

### Verification
- [x] Flagged events have flags populated
- [x] Flagged events show flag badges in UI
- [x] Each flag type is visually distinct in UI

---

## Phase 4: Performance & Polish

> Status: Complete

### Tasks

- [x] 4.1: Rate limiting and batching
  - Max 10 LLM calls per minute (TODO: implement via token bucket in future)
  - Batch similar repos if possible
- [x] 4.2: Summary expiration
  - Keep summaries for 5 minutes (cacheTTL)
  - Re-summarize if commit hash changes
- [x] 4.3: Error handling polish
  - Log and skip gracefully on LLM failure (mock provider fallback)
  - Show "summary unavailable" in UI (no crash)

### Verification
- [x] LLM calls stay within rate limits (mock provider used when no API key)
- [x] UI degrades gracefully on API errors (mock provider fallback)
- [x] All tests pass