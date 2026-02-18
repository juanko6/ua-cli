# Data Model: Schedule Module (Updated)

**Feature**: 001-devops-schedule-polish

## Entities

### Event (updated)

| Field    | Type        | Description                        |
|----------|-------------|------------------------------------|
| ID       | string      | Optional UACloud identifier        |
| Title    | string      | Subject name + code                |
| Start    | time.Time   | Start datetime                     |
| End      | time.Time   | End datetime                       |
| Location | string      | Room/building (if available)       |
| Type     | EventType   | Teoria, Practica, Seminario, etc.  |

### EventType (enum)

| Value        | Color Mapping                |
|--------------|------------------------------|
| Teoria       | TBD (green tones)            |
| Practica     | TBD (blue tones)             |
| Seminario    | TBD (orange tones)           |
| Desconocido  | Fallback for unknown colors  |

### Dedup Key

Composite of `Title + Start.Unix() + End.Unix()`. Used as map key in service layer.

## Validation Rules

- `Start` MUST be before `End`
- `Title` MUST NOT be empty
- Events with `Start` outside the target week MUST be filtered out
