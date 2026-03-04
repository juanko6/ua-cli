# Data Model: ua now Command

**Feature**: 002-now-command | **Date**: 2026-03-04

No new domain entities are being created for this feature. The command relies 100% on the existing `schedule.Event` data model established in Phase 1.

## Entities (Existing)

### `schedule.Event`

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique event ID from UACloud |
| Title | string | Subject Name (e.g., ESTRUCTURA DE LOS COMPUTADORES) |
| Start | time.Time | Parsed start timestamp of the class |
| End | time.Time | Parsed end timestamp of the class |
| Type | enum | Theory, Practice, Seminar, Unknown |
| Location | string | Classroom Identifier (e.g., A3/0007) |

## State Transitions (CLI Logic)

The command relies on comparing the entity's `Start` and `End` properties with `time.Now()`.

```text
time.Now() < Start:
-> Event is UPCOMING. Status = "PRÓXIMA"

Start <= time.Now() < End (and End - Now > 30m):
-> Event is ONGOING. Status = "AHORA MISMO"

Start <= time.Now() < End (and End - Now <= 30m):
-> Event is ENDING SOON. Status = "PRÓXIMA" (Shifts reference to chronologically next Event).
```
