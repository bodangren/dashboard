# Track: Commit Search Query Deduplication

## Problem
`CommitSearchQuery`, `parseRelativeDate`, and `relativeDateRegex` are duplicated in both `internal/api/handlers.go` and `internal/search/query.go`. The API layer should import from the search package instead of redefining these types.

## Goal
Eliminate duplication by moving canonical definitions to `internal/search` and having `internal/api` import them.

## Acceptance Criteria
- [ ] `CommitSearchQuery`, `parseRelativeDate`, and `relativeDateRegex` live only in `internal/search`
- [ ] `internal/api/handlers.go` imports from `internal/search` with zero duplication
- [ ] All Go tests pass (`go test ./...`)
- [ ] Build succeeds (`go build ./...`)
- [ ] No functional regressions in commit search behavior

## Related Tech Debt
- DUP-01: CommitSearchQuery duplicated in api and search packages
