# Plan: Commit Search Query Deduplication

## Phase 1: Extract and Export (TDD)
- [ ] Write tests verifying search package exports CommitSearchQuery and parseRelativeDate
- [ ] Ensure relativeDateRegex is package-internal or exported as needed
- [ ] Tests pass

## Phase 2: Remove API Duplication
- [ ] Delete duplicate definitions from internal/api/handlers.go
- [ ] Add import of internal/search types
- [ ] Update call sites to use search.CommitSearchQuery etc.
- [ ] Tests pass

## Phase 3: Verification
- [ ] Run full Go test suite
- [ ] Run go build ./...
- [ ] Manual smoke test of commit search UI
- [ ] Update tech-debt.md to mark DUP-01 resolved
- [ ] Commit and push
