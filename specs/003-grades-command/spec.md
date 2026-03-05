# Feature Specification: `ua grades` Command

**Feature Branch**: `003-grades-command`  
**Created**: 2026-03-05  
**Status**: Draft  
**Input**: User description: "quiero agregar un comando que me muestre mis notas actuales de todas las asignaturas, detectando cuando hay nuevas notas publicadas y mostrando un indicador visual"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Grades Overview (Priority: P1)

As a university student tracking my academic performance, I want to type `ua grades` and see a complete overview of my current grades across all subjects, so I can quickly assess my standing without logging into the web portal.

**Why this priority**: Grades are the most critical academic information after schedule. Students check grades frequently to track progress, identify struggling subjects, and anticipate final results.

**Independent Test**: Can be fully tested by running `ua grades` and verifying all enrolled subjects are displayed with their current grades, averages, and status indicators.

**Acceptance Scenarios**:

1. **Given** the user has enrolled subjects with grades, **When** they run `ua grades`, **Then** the CLI displays a table with subject name, current grade, average, and status (approved/pending).
2. **Given** some subjects have no grades yet, **When** they run `ua grades`, **Then** those subjects show "No calificación" or equivalent.
3. **Given** the user runs `ua grades --json`, **Then** the system outputs structured JSON data compatible with external tools.

---

### User Story 2 - New Grades Detection (Priority: P2)

As a busy student, I want the CLI to automatically detect and highlight newly published grades since my last check, so I can immediately identify what needs my attention.

**Why this priority**: Provides immediate value by reducing the need for manual checking. Students want to know when new grades are available without constantly checking.

**Independent Test**: Can be tested by running `ua grades` twice and verifying that newly added grades are highlighted on the second run (requires mocking or real grade changes).

**Acceptance Scenarios**:

1. **Given** the user hasn't checked grades since last week, **When** they run `ua grades`, **Then** any newly published grades are highlighted with a "🆕" indicator or similar visual cue.
2. **Given** there are no new grades since last check, **When** they run `ua grades`, **Then** the system displays "✅ No hay nuevas calificaciones" or equivalent.
3. **Given** new grades are detected, **When** the user runs `ua grades`, **Then** the system updates the internal tracking of last checked date.

---

### Edge Cases

- What happens when UACloud returns an empty grades response? → Should display a helpful message about no grades available.
- What happens when there's a network error fetching grades? → Should gracefully handle and suggest trying again.
- How does the system handle subjects with partial grading (multiple assessments)? → Should show current average and indicate incomplete grading.
- What happens when grade format changes (e.g., from numbers to letters)? → Should be resilient to format changes through parsing abstraction.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST fetch current grades for all enrolled subjects from UACloud.
- **FR-002**: System MUST display grades in a clear tabular format with columns: Subject, Grade, Average, Status.
- **FR-003**: System MUST detect and highlight newly published grades since the last command execution.
- **FR-004**: System MUST provide JSON output via `--json` flag for automation purposes.
- **FR-005**: System MUST track last execution date to detect new grades.
- **FR-006**: System MUST handle subjects with no grades gracefully.
- **FR-007**: System MUST calculate and display current average across all subjects.
- **FR-008**: System MUST indicate approval status (approved/pending/requiring attention).
- **FR-009**: System MUST provide filtering options (e.g., `--approved`, `--pending`).

### Key Entities

- **Grade Record**: Contains subject name, current grade, assessment count, average, and approval status.
- **Grade Change**: Represents a newly published grade since last check, with timestamp and subject identifier.
- **Grade Tracker**: Internal component that tracks last check date and detected changes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can view their complete grades overview in under 0.8 seconds.
- **SC-002**: 100% accuracy in detecting and highlighting new grades on subsequent runs.
- **SC-003**: Support for all grade formats used by UACloud (numerical, percentage, letter grades).
- **SC-004**: JSON output must be valid and contain all grade data for external tool integration.