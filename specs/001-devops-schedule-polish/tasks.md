# Tasks: DevOps Infrastructure & Schedule Polish

**Input**: Design documents from `/specs/001-devops-schedule-polish/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, quickstart.md ✅

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1–US5)
- Exact file paths included in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Rename branch, configure project metadata, baseline tag.

- [x] T001 Rename default branch from `master` to `main` via `gh repo rename-branch` and update local tracking
- [x] T002 Create `CHANGELOG.md` at repository root with Keep a Changelog format and `[0.1.0]` entry

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: CI/CD pipelines must exist before feature branches can use PR workflow.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T003 Create `.github/workflows/ci.yml` with lint (golangci-lint), build, and test jobs triggered on push and PR
- [x] T004 [P] Create `.goreleaser.yml` with cross-compile targets (linux/amd64, darwin/arm64, windows/amd64)
- [x] T005 [P] Create `.github/workflows/release.yml` triggered on `v*` tags, running GoReleaser
- [ ] T006 Tag current commit as `v0.1.0` and push tag to trigger first release

**Checkpoint**: CI runs on push. Tagging triggers release. PR workflow is functional.

---

## Phase 3: User Story 1 — Automated Quality Gates (Priority: P1) 🎯 MVP

**Goal**: Every push/PR triggers CI (lint + build + test) and blocks merge on failure.

**Independent Test**: Push a broken commit → CI must fail → PR blocked.

### Implementation for User Story 1

- [x] T007 [US1] Verify `ci.yml` runs on a test branch push by creating `fix/ci-test` branch and pushing a dummy commit
- [x] T008 [US1] Verify PR checks block merge by opening a test PR with a failing build

**Checkpoint**: CI pipeline validated end-to-end.

---

## Phase 4: User Story 2 — Automated Releases (Priority: P1)

**Goal**: Tagging `v*` builds cross-platform binaries and publishes a GitHub Release.

**Independent Test**: Push `v0.1.0` tag → GitHub Release with 3 binaries.

### Implementation for User Story 2

- [x] T009 [US2] Verify GoReleaser produces binaries locally via `goreleaser release --snapshot --clean`
- [x] T010 [US2] Verify release pipeline by confirming `v0.1.0` tag creates a GitHub Release with assets

**Checkpoint**: Release pipeline validated.

---

## Phase 5: User Story 3 — Accurate Weekly Schedule (Priority: P1)

**Goal**: `ua schedule` shows only events within the current Monday–Sunday range.

**Independent Test**: Run `ua schedule --json | jq length` and verify ≤ 30 events, all within the current week.

### Implementation for User Story 3

- [x] T011 [US3] Add week-filtering logic in `internal/service/schedule/service.go` — filter events where `Start` falls within Monday 00:00 to Sunday 23:59 of the target week
- [x] T012 [US3] Strip trailing newline from cookie file read in `cmd/ua-cli/schedule.go`
- [x] T013 [US3] Add `--next` flag to `cmd/ua-cli/schedule.go` to show next week's events

**Checkpoint**: `ua schedule` returns only current week events.

---

## Phase 6: User Story 4 — Classified Event Types (Priority: P2)

**Goal**: Events display their type (Teoría, Práctica, Seminario) based on API color data.

**Independent Test**: Run `ua schedule --json` and verify ≥ 80% of events have `type != "Desconocido"`.

### Implementation for User Story 4

- [x] T014 [US4] Log unique color values from real API response by adding debug output in `internal/adapters/uacloud/parser.go`
- [x] T015 [US4] Update `classifyEvent()` in `internal/adapters/uacloud/parser.go` with real color-to-type mapping
- [x] T016 [US4] Remove debug logging added in T014

**Checkpoint**: Events display correct types.

---

## Phase 7: User Story 5 — No Duplicate Events (Priority: P2)

**Goal**: Each unique event (by title + start + end) appears only once.

**Independent Test**: Run `ua schedule --json | jq '[.[] | {t:.title,s:.start,e:.end}] | unique | length'` equals total length.

### Implementation for User Story 5

- [x] T017 [P] [US5] Add `DeduplicateEvents()` function in `internal/service/schedule/service.go` using composite key `Title+Start.Unix()+End.Unix()`
- [x] T018 [US5] Call `DeduplicateEvents()` in `GetScheduleForWeek()` before returning events

**Checkpoint**: No duplicate events in output.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Cleanup, documentation, and final validation.

- [x] T019 [P] Remove all DEBUG print statements from `internal/adapters/uacloud/adapter.go`
- [x] T020 [P] Update `CHANGELOG.md` with `[0.2.0]` entry covering all changes
- [x] T021 [P] Update `specs/001-devops-schedule-polish/spec.md` status from Draft → Complete
- [ ] T022 Tag `v0.2.0` and push to trigger release with all fixes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Phase 2 (validates CI)
- **US2 (Phase 4)**: Depends on Phase 2 (validates releases)
- **US3 (Phase 5)**: Depends on Phase 2 — independent of US1/US2
- **US4 (Phase 6)**: Depends on Phase 2 — independent of other stories
- **US5 (Phase 7)**: Depends on Phase 2 — independent of other stories
- **Polish (Phase 8)**: Depends on all desired stories being complete

### User Story Dependencies

- **US1**: Can start after Phase 2 — No cross-story dependency
- **US2**: Can start after Phase 2 — No cross-story dependency
- **US3**: Can start after Phase 2 — No cross-story dependency
- **US4**: Can start after Phase 2 — No cross-story dependency
- **US5**: Can start after Phase 2 — No cross-story dependency

### Parallel Opportunities

- T003, T004, T005 can run in parallel (different files)
- US3, US4, US5 can all proceed in parallel after Phase 2
- T017 is parallelizable with T015/T016
- T019, T020, T021 can all run in parallel

---

## Parallel Example: User Story 5

```bash
# These can run in parallel (different functions, different concerns):
Task T017: "Add DeduplicateEvents() in service.go"
Task T015: "Update classifyEvent() in parser.go"
```

---

## Implementation Strategy

### MVP First (US1 + US2 + US3)

1. Complete Phase 1: Setup (rename branch, CHANGELOG)
2. Complete Phase 2: CI/CD pipelines
3. Complete Phase 3: Validate CI
4. Complete Phase 4: Validate releases
5. Complete Phase 5: Week filter (most impactful bug fix)
6. **STOP and VALIDATE**: `ua schedule` shows current week only
7. Tag `v0.2.0`

### Incremental Delivery

1. Setup + CI/CD → foundation ready
2. Add week filter → schedule is usable (MVP!)
3. Add event types → schedule is informative
4. Add dedup → schedule is clean
5. Polish → release v0.2.0

---

## Notes

- [P] tasks = different files, no dependencies
- No test tasks generated (tests not requested in spec)
- Each user story is independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
