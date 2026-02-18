# Feature Specification: DevOps Infrastructure & Schedule Polish

**Feature Branch**: `001-devops-schedule-polish`  
**Created**: 2026-02-17  
**Status**: Draft  
**Input**: User description: "Configurar infraestructura DevOps (GitHub Flow, SemVer, CI/CD con GitHub Actions + GoReleaser) y polish del módulo Schedule (filtro semanal, colores de eventos, deduplicación)"

## User Scenarios & Testing

### User Story 1 - Automated Quality Gates (Priority: P1)

As a developer, I want every push and pull request to be automatically validated (lint, build, test) so that broken code never reaches the main branch.

**Why this priority**: Without CI, bugs can be merged silently and compound.

**Independent Test**: Push a commit with a syntax error; the CI pipeline should fail and block merge.

**Acceptance Scenarios**:

1. **Given** a push to any branch, **When** CI runs, **Then** linting, building, and testing complete within 5 minutes.
2. **Given** a PR to `main`, **When** CI fails, **Then** the PR is blocked from merging.

---

### User Story 2 - Automated Releases (Priority: P1)

As a maintainer, I want tagging a version (e.g. `v0.2.0`) to automatically build and publish cross-platform binaries on GitHub Releases.

**Why this priority**: Manual releases are error-prone and time-consuming.

**Independent Test**: Push a `v0.1.0` tag; a GitHub Release should be created with binaries for Linux, macOS, and Windows.

**Acceptance Scenarios**:

1. **Given** a tag `v*` is pushed, **When** the release pipeline runs, **Then** binaries for linux/amd64, darwin/arm64, windows/amd64 are published.
2. **Given** the release pipeline completes, **Then** a CHANGELOG entry for the version is available.

---

### User Story 3 - Accurate Weekly Schedule (Priority: P1)

As a student, I want `ua schedule` to show only classes for the current week so that the output is useful and not cluttered with the entire semester.

**Why this priority**: Currently the schedule shows all events from the semester, making it unusable.

**Independent Test**: Run `ua schedule` on a Monday; only events from that Monday to Sunday should appear.

**Acceptance Scenarios**:

1. **Given** events spanning the entire semester, **When** `ua schedule` runs, **Then** only events within the current week (Mon-Sun) are displayed.
2. **Given** `--next` flag, **When** run, **Then** show next week's events only.

---

### User Story 4 - Classified Event Types (Priority: P2)

As a student, I want each class to show its type (Teoría, Práctica, Seminario) so I can distinguish between session types.

**Why this priority**: All events currently show "Desconocido", reducing the value of the information.

**Independent Test**: Run `ua schedule --json`; events should have their `type` field populated based on the API color data.

**Acceptance Scenarios**:

1. **Given** an event with a known color code, **When** displayed, **Then** it shows the correct type label.

---

### User Story 5 - No Duplicate Events (Priority: P2)

As a student, I want each class to appear only once per time slot so the schedule is clean and readable.

**Why this priority**: Duplicates clutter the output and confuse users.

**Independent Test**: Run `ua schedule`; no two rows should have identical day, time, and subject.

**Acceptance Scenarios**:

1. **Given** raw API data with duplicates, **When** processed, **Then** only one instance of each unique event (by title+start+end) remains.

---

### Edge Cases

- What happens when the API returns an empty array? → Show "No classes scheduled this week."
- What happens when the cookie has expired? → Show a clear error: "Session expired. Please update your cookie."
- What happens when the user has no internet? → If cache exists, use cached data; otherwise, show "No internet and no cached data."

## Requirements

### Functional Requirements

- **FR-001**: The project MUST use GitHub Flow (all changes via feature branches → PRs → main).
- **FR-002**: The CI pipeline MUST run lint, build, and test on every push and PR.
- **FR-003**: The release pipeline MUST build cross-platform binaries when a `v*` tag is pushed.
- **FR-004**: The project MUST follow Semantic Versioning (MAJOR.MINOR.PATCH).
- **FR-005**: `ua schedule` MUST filter events to the requested week only.
- **FR-006**: `ua schedule` MUST classify events by type based on available API metadata (color).
- **FR-007**: `ua schedule` MUST deduplicate events with identical title, start time, and end time.

### Key Entities

- **CI Pipeline**: A set of automated jobs (lint, build, test) triggered on code changes.
- **Release Pipeline**: An automated workflow that produces distributable binaries from tagged commits.
- **Schedule Event**: An academic activity with title, start/end time, location, and type.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every PR to `main` is validated by CI in under 5 minutes.
- **SC-002**: Tagged releases produce downloadable binaries within 10 minutes.
- **SC-003**: `ua schedule` shows ≤ 30 events per week (no semester-wide dump).
- **SC-004**: 0 duplicate events appear in the output for any given week.
- **SC-005**: ≥ 80% of events have a classified type (not "Desconocido").

## Assumptions

- GitHub Actions is the CI/CD platform (free tier sufficient for this project).
- GoReleaser handles cross-compilation and release artifact generation.
- The color codes in UACloud's JSON response are consistent and can be mapped to event types.
- The cookie-based authentication remains stable for the duration of MVP development.
