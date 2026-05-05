# Plan: Multi-Platform Support & Advanced Git Features

## Context

Track: platform-git-features_20260505

This track implements branch management and stash viewing to enhance the git dashboard beyond simple pull and log operations. Multi-user authentication is explicitly out of scope per the product definition (solo developer use).

## Tasks

### Phase 1: Branch Management API

#### Task 1.1: Add branch list and checkout to git package
- [x] Add `GetBranches(repoPath string) ([]string, error)` in internal/git/branch.go
- [x] Add `GetCurrentBranch(repoPath string) (string, error)`
- [x] Add `CreateBranch(repoPath, branchName string) error`
- [x] Add `DeleteBranch(repoPath, branchName string) error`
- [x] Add `CheckoutBranch(repoPath, branchName string) error`
- [x] Write unit tests for all branch functions

#### Task 1.2: Add branch handlers to API
- [x] Add `GetBranchesFunc func(repoPath string) ([]string, error)` to HandlerConfig
- [x] Add `CheckoutBranchFunc func(repoPath, branch string) error` to HandlerConfig
- [x] Add `CreateBranchFunc func(repoPath, branch string) error` to HandlerConfig
- [x] Add `DeleteBranchFunc func(repoPath, branch string) error` to HandlerConfig
- [x] Add GET /api/branches endpoint
- [x] Add POST /api/branches/checkout endpoint
- [x] Add POST /api/branches (create)
- [x] Add DELETE /api/branches endpoint
- [x] Write tests for branch handlers

### Phase 2: Stash Management API

#### Task 2.1: Add stash functions to git package
- [x] Add `GetStashList(repoPath string) ([]StashEntry, error)` in internal/git/stash.go
- [x] Add `ApplyStash(repoPath string, index int) error`
- [x] Add `DropStash(repoPath string, index int) error`
- [x] Write unit tests for stash functions

#### Task 2.2: Add stash handlers to API
- [x] Add `GetStashFunc func(repoPath string) ([]StashEntry, error)` to HandlerConfig
- [x] Add `ApplyStashFunc func(repoPath string, index int) error` to HandlerConfig
- [x] Add `DropStashFunc func(repoPath string, index int) error` to HandlerConfig
- [x] Add GET /api/stash endpoint
- [x] Add POST /api/stash/apply endpoint
- [x] Add DELETE /api/stash/{index} endpoint
- [x] Write tests for stash handlers

### Phase 3: Frontend Integration

#### Task 3.1: Branch picker UI
- [x] Add branch badge to project card showing current branch
- [x] Add branch dropdown with list of local branches
- [x] Add "Create branch" option in dropdown
- [x] Add "Delete branch" option in dropdown (with confirmation)
- [x] Wire up POST /api/branches/checkout on branch selection

#### Task 3.2: Stash viewer UI
- [x] Add stash toggle button to project card
- [x] Create stash panel overlay showing stash list
- [x] Add "Apply" and "Drop" buttons per stash entry
- [x] Wire up stash API endpoints

#### Task 3.3: Verification
- [x] Run `go test ./...`
- [x] Run `go build ./...`
- [x] Start dev server, verify in browser via browser-harness
- [x] Check for console errors
- [x] Commit