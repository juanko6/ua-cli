# Tasks: Schedule Module

**Input**: [plan.md](./plan.md)

## Phase 1: Core Domain (Hexagonal)

- [x] **T001** [Story: US1] Define `Event` struct in `internal/domain/schedule/entity.go`.
- [x] **T002** [Story: US1] Define `Repository` interface in `internal/domain/schedule/repository.go`.
- [x] **T003** [Story: US1] Implement `ScheduleService` (Business Logic: Caching rules 24h, Date calc).

## Phase 2: Adapters (Infrastructure)

- [x] **T004** [Story: US1] Implement `UACloudAdapter` (HTTP Client).
- [x] **T005** [Story: US1] Implement `JSONFileRepo` in `internal/adapters/repo/json.go` (Load/Save/Stat).
- [x] **T006** [Story: US1] Implement `UAParser` (Resilient HTML/JSON parsing).

## Phase 3: CLI & Adaptive UI

- [x] **T007** [Story: US2] Implement `TUI Presenter` (Bubbletea table) in `internal/adapters/presenter/tui.go`.
- [x] **T008** [Story: US3] Implement `Text Presenter` (Tabwriter) in `internal/adapters/presenter/text.go`.
- [x] **T009** [Story: US1] Create `cmd/ua-cli/schedule.go` with TTY detection logic.
- [x] **T010** [Story: US2] Wire up navigation flags (`--next`, `--prev`) to Service.
- [x] **T011** [Story: US3] Wire up `--json` flag.

## Phase 4: Polish

- [ ] **T012** [Story: US1] Add integration tests (Mock Server -> Adapter -> Service).
- [ ] **T013** [Polish] Verify 24h cache expiry logic works offline.
