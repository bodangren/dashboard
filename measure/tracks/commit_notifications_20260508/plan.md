# Plan: Commit and Agent Failure Notifications

## Phase 1: Data Model & Storage (TDD)
- [ ] Write migration and tests for `notifications` table (id, type, message, repo, read, created_at)
- [ ] Write tests for NotificationStore CRUD
- [ ] Implement NotificationStore

## Phase 2: Rules Engine (TDD)
- [ ] Write tests for rule evaluation against commit/agent events
- [ ] Implement RuleEngine with configurable thresholds
- [ ] Wire RuleEngine into Hub event pipeline

## Phase 3: WebSocket Delivery (TDD)
- [ ] Write tests for WebSocket notification broadcast
- [ ] Add `notification` message type to Hub protocol
- [ ] Broadcast on new notification creation

## Phase 4: Frontend UI (TDD)
- [ ] Write tests for NotificationBell component
- [ ] Implement NotificationBell with unread count
- [ ] Write tests for NotificationPanel component
- [ ] Implement NotificationPanel with mark-read and dismiss

## Phase 5: Integration & Verification
- [ ] End-to-end test: commit triggers notification → appears in UI
- [ ] End-to-end test: agent failure triggers notification → appears in UI
- [ ] Update tracks.md and commit
