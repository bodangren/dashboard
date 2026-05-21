# Plan: Resolve API-Search Import Cycle and CommitSearchQuery Deduplication

## Phase 1: Diagnose and Break Import Cycle
- [ ] Write test proving the current import cycle exists
- [ ] Identify why search depends on git (or api)
- [ ] Extract shared types to a new `internal/search/types.go` or `internal/gittypes` package that has no upstream dependencies
- [ ] Verify cycle is broken with `go build ./...`

## Phase 2: Migrate CommitSearchQuery to Search Package
- [ ] Move `CommitSearchQuery` struct and validation to `internal/search/query.go`
- [ ] Write tests for `CommitSearchQuery` in search package
- [ ] Update `internal/api/handlers.go` to import from search
- [ ] Remove duplicate from api package
- [ ] Verify handlers tests pass

## Phase 3: Migrate Date Parsing Helpers
- [ ] Move `parseRelativeDate` and `relativeDateRegex` to search package
- [ ] Write tests for relative date parsing
- [ ] Update api handlers to use search package versions
- [ ] Remove duplicates from api package
- [ ] Verify integration

## Phase 4: Final Verification
- [ ] Run full test suite (`go test ./...`)
- [ ] Run build (`go build ./...`)
- [ ] Update tech-debt.md to mark DUP-01 Resolved
- [ ] Commit and push
