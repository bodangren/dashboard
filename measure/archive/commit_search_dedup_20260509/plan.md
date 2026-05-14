# Plan: Commit Search Query Deduplication

## Phase 1: Analysis (COMPLETED)
- [x] Investigate duplicate definitions in api and search packages
- [x] Identify import cycle issue: git->api->search->git prevents straightforward dedup

## Phase 2: Findings
The proposed refactor (moving CommitSearchQuery to search package only) is blocked by a pre-existing import cycle:
- git/log.go imports api (ToAPICommit method)
- If api imports search to use search.CommitSearchQuery, it creates: search->git->api->search cycle

The existing duplication in handlers.go (lines 526-582) actually exists to break this cycle. The plan as specified cannot be implemented without first resolving the git->api dependency.

## Phase 3: Conclusion
- [x] DUP-01 remains Open due to architectural constraint
- [x] Pre-existing import cycle must be resolved before dedup can proceed
- [x] tech-debt.md notes this as DUP-01 Open

## Note
The duplication is localized to handlers.go and does not cause runtime issues.
Both CommitSearchQuery implementations are functionally equivalent.
