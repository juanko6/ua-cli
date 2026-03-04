# Feature Specification: `ua now` Command

**Feature Branch**: `002-now-command`  
**Created**: 2026-03-04  
**Status**: Draft  
**Input**: User description: "quiero agregar a schedule algo, que si quiero saber que clase tengo que solo me traiga la clase mas cercana a mi hora actual, es decir si tengo clase de 10 a 12 de estructura y son las 9:50 o si esta dentro de la hora de la clase hasta faltando 30min para acabarse, pues que me diga tienes nombre asignatura, la localidad, hora inicia, hora fin, ya si falta menos de 30 min muestra la siguiente clase. me gusta mas el now, es mas corta vamos con ello"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Instant Next Class Lookup (Priority: P1)

As a university student managing a busy schedule, I want to type a single, short command to instantly see my current or immediate next class, so I know exactly where I need to be without parsing through my entire weekly schedule.

**Why this priority**: Directly targets the project's North Star metric of minimizing the time it takes to get critical academic information. Looking up "what class do I have right now" is the most frequent use case for a schedule.

**Independent Test**: Can be fully tested by running `ua now` at various times of the day (before first class, during a class, right before a class ends) and verifying the output matches the expected temporal rules.

**Acceptance Scenarios**:

1. **Given** the current time is before the first class of the day, **When** the user runs `ua now`, **Then** the CLI displays details of that first class.
2. **Given** the user is currently in a class and there are 45 minutes remaining, **When** they run `ua now`, **Then** the CLI displays the details of the ongoing class.
3. **Given** the user is currently in a class and there are 15 minutes remaining, **When** they run `ua now`, **Then** the CLI displays the details of the *next* chronological class of the day.
4. **Given** the user has finished all classes for the day, **When** they run `ua now`, **Then** the CLI displays a friendly message indicating there are no more classes.

---

### Edge Cases

- What happens when the user runs the command on a day with no scheduled classes (e.g., weekends)? -> Should gracefully inform the user they have a free day.
- What happens if there is less than 30 minutes remaining in the last class of the day? -> Should indicate there are no *more* classes (i.e., the day is effectively over).
- How does the system handle time zone shifts when comparing the class times with the local system time? -> Must ensure both times are comparable within the local timezone.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST identify the current system time and filter the student's weekly schedule to isolate only events occurring on the current calendar day.
- **FR-002**: System MUST display the details of an upcoming class if the current time is before that class's start time and no other classes precede it.
- **FR-003**: System MUST display the details of an ongoing class if the current time falls within the class duration AND there is strictly more than 30 minutes remaining until the class ends.
- **FR-004**: System MUST automatically advance to displaying the *next* chronologically scheduled class of the day if the current ongoing class has 30 minutes or less remaining.
- **FR-005**: System MUST display a clear "no more classes" or "free day" message if all classes for the day have concluded or none exist.
- **FR-006**: System MUST output the following data points for the identified class: subject name, location/classroom, start time, end time, and class type (theory/practice).
- **FR-007**: System MUST provide output in both terminal-friendly UI and raw JSON format (via a `--json` flag).

### Key Entities

- **Schedule Event**: Represents a class block. Contains temporal bounds (start, end) and metadata (title, location, type). Logic depends strictly on the temporal bounds relative to "now".

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can execute the command and view their immediate next class data in under 0.5 seconds.
- **SC-002**: 100% adherence to the 30-minute transition rule: tests must confirm the transition from "current class" to "next class" occurs exactly at the T-minus 30 minute mark.
- **SC-003**: Reduction in the mechanical strokes needed: replacing `ua schedule | grep <today>` with `ua now`.
