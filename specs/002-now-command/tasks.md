# Tasks: ua now Command

**Input**: Design documents from `/specs/002-now-command/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, quickstart.md ✅

**Tests**: Not explicitly requested in spec.
**Organization**: Tasks are grouped by user story for independent implementation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1)
- Exact file paths included in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- *(No new setup tasks required; reusing existing Go modules and CLI structure)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- *(No new foundational tasks required; reusing existing `ScheduleService` and caching layer)*

---

## Phase 3: User Story 1 - Instant Next Class Lookup (Priority: P1) 🎯 MVP

**Goal**: Type a single command to see the current or immediate next class based on a 30-minute threshold.

**Independent Test**: Can be fully tested by running `ua now` and validating the output matches the temporal rules.

### Implementation for User Story 1

- [x] T001 [US1] Create basic Cobra command structure in `cmd/ua-cli/now.go` with service initialization
- [x] T002 [US1] Implement schedule filtering and 30-minute threshold sorting logic in `cmd/ua-cli/now.go`
- [x] T003 [P] [US1] Implement Lipgloss card UI rendering (`renderNowCard`) in `cmd/ua-cli/now.go`
- [x] T004 [P] [US1] Implement JSON output support (`--json`) for automated scripts in `cmd/ua-cli/now.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T005 [P] Run `go fmt ./...` and `go test ./...` to ensure code is clean and passes all existing tests
- [x] T006 Verify the feature exactly matches the scenarios in `specs/002-now-command/quickstart.md`
- [x] T007 Commit all changes to the `002-now-command` branch

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup / Foundational**: Skipped (existing infrastructure is fully sufficient)
- **User Stories (Phase 3)**: Unblocked immediately
- **Polish (Phase 4)**: Depends on User Story 1 completion

### Parallel Opportunities

- T003 and T004 can be implemented in parallel after T002 (logic routing complete)
- T005 and T006 are parallel validation steps
