# Feature Specification: Schedule Module (Horarios)

**Feature Branch**: `feature-schedule`
**Created**: 2026-02-17
**Status**: Draft
**Input**: User description: "implementar el módulo de Horarios"

## User Scenarios & Testing

### User Story 1 - View Weekly Schedule (Priority: P1)
Current User wants to see their classes for the current week to know where to go next.

**Why this priority**: Core value proposition. Reduces friction of logging into web portal.

**Independent Test**:
- Run `ua schedule`.
- Verify output shows Monday-Sunday events for the current week.
- Verify <1s execution time.

**Acceptance Scenarios**:
1. **Given** a valid session, **When** user runs `ua schedule`, **Then** show ASCII table with events sorted by day/time.
2. **Given** expired session, **When** user runs `ua schedule`, **Then** prompt to login or update cookie.

---

### User Story 2 - Navigation (Priority: P1)
User wants to check next week's sorting or previous week's attendance.

**Acceptance Scenarios**:
1. **Given** a valid session, **When** user runs `ua schedule --next` (or `-w 1`), **Then** show schedule for next week.
2. **Given** a valid session, **When** user runs `ua schedule --prev` (or `-w -1`), **Then** show schedule for previous week.

---

### User Story 3 - Machine Readable Output (Priority: P2)
User wants to pipe the schedule to another tool (e.g., Rainmeter, script).

**Acceptance Scenarios**:
1. **Given** valid data, **When** user runs `ua schedule --json`, **Then** output raw JSON data without ASCII decoration.

---

## Requirements

### Functional Requirements

- **FR-001**: System MUST fetch schedule events from `cvnet.cpd.ua.es`.
- **FR-002**: System MUST parse the proprietary JSON/HTML format from UA.
- **FR-003**: System MUST group events by Day (Monday-Sunday).
- **FR-004**: System MUST display: Title, Start Time, End Time, Location (Aula), and Type (Teor/Prac).
- **FR-005**: System MUST cache the schedule locally in a JSON file (`~/.ua-cli/cache/schedule.json`) to enable offline viewing.
- **FR-006**: System MUST respect "Offline First": try cache first, fallback to network only if cache expired (> 24 hours) or forced refresh.

## Clarifications
### Session 2026-02-17
- Q: Cache Invalidation Strategy → A: 24 Hours (Balanced).
- Q: Persistence Engine → A: JSON File (Simple).
- Q: UI Output Strategy → A: Adaptive (Bubbletea if TTY, Text if Pipe).

### Key Entities

- **Event**:
    - `ReferenceId` (string)
    - `Title` (string)
    - `Start` (time)
    - `End` (time)
    - `Location` (string)
    - `EventType` (enum: Theory, Practice, Seminar)

## Success Criteria

### Measurable Outcomes

- **SC-001**: `ua schedule` renders in under **500ms** (warm cache).
- **SC-002**: `ua schedule` renders in under **2s** (cold network fetch).
- **SC-003**: ZERO "panic" errors on malformed API responses (robust parsing).
- **SC-004**: Works 100% offline if data was fetched previously.

### Edge Cases
- **Empty Week**: Show "No classes" message instead of empty table.
- **Holidays**: Handle days without events gracefully.
- **API Change**: Detect schema mismatch and warn user "UA API changed, update CLI".
