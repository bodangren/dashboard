# Spec: Resolve API-Search Import Cycle and CommitSearchQuery Deduplication

## Problem
DUP-01 documents that `CommitSearchQuery`, `parseRelativeDate`, and `relativeDateRegex` are duplicated across `internal/api/handlers.go` and `internal/search/query.go`. The API layer should import from the search package, but an import cycle prevents this: `git -> api -> search -> git`.

## Goal
Break the import cycle so that search package types can be imported by api without creating a cycle. Then deduplicate CommitSearchQuery and related helpers.

## Acceptance Criteria
- [ ] No import cycle between api and search packages
- [ ] `CommitSearchQuery` lives in only one package (search)
- [ ] `parseRelativeDate` and `relativeDateRegex` live in only one package (search)
- [ ] All existing tests pass after refactoring
- [ ] No behavioral changes to commit search functionality
