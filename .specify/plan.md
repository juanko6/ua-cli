# Implementation Plan: Schedule Module

**Branch**: `feature-schedule` | **Date**: 2026-02-17 | **Spec**: [spec.md](./spec.md)

## Summary
Implement `ua schedule` to fetch, cache, and display the student's weekly schedule.
**Key constraint**: Must be "Adaptive" — beautiful TUI (Bubbletea) for humans, simple text/JSON for scripts (`|`).

## Technical Context
- **Language**: Go 1.22
- **Architecture**: Hexagonal (Ports & Adapters)
- **External API**: `cvnet.cpd.ua.es`
- **Persistence**: Flat JSON file (`~/.ua-cli/cache/schedule.json`). No SQLite.
- **UI**: Bubbletea (TUI) + `text/tabwriter` (CLI).

## Constitution Check
- **Speed**: JSON Cache (<10ms read).
- **UX**: "Unhinged" Bubbletea view for interactive sessions. Standard text for pipes.
- **Privacy**: Local-first storage.

## Project Structure

```text
ua-cli/
├── internal/
│   ├── domain/
│   │   └── schedule/       # Entities (Event) + Interface (Repository)
│   ├── adapters/
│   │   ├── uacloud/        # HTTP Client (UA API)
│   │   ├── repo/           # JSON Filesystem Repo
│   │   └── presenter/      # UI Logic (Adaptive)
│   └── service/            # Application Logic (GetWeeklySchedule)
├── cmd/
│   └── ua-cli/
│       └── schedule.go     # CLI Command (Wires everything)
```

## Technical Design

### Data Model
```go
type Event struct {
    ID       string
    Title    string
    Start    time.Time
    End      time.Time
    Location string
    Type     string // Theory, Practice
}
```

### Persistence (JSON)
- File: `~/.ua-cli/cache/schedule.json`
- Format: `{"updated_at": "...", "events": [...]}`
- TTL: 24 hours (configurable).

### Adaptive UI Strategy
Check `isatty` (via `golang.org/x/term` or `muesli/termenv`).
- **If TTY**: Launch Bubbletea Model (Interactive table, navigation with arrows).
- **If !TTY**: Print Table via `text/tabwriter` (Clean, grep-friendly).
- **If --json**: Print Raw JSON.
