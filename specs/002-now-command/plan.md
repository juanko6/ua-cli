# Implementation Plan: ua now Command

**Branch**: `002-now-command` | **Date**: 2026-03-04 | **Spec**: [spec.md](file:///c:/Users/juankodev/Desktop/ua-cli/specs/002-now-command/spec.md)
**Input**: Feature specification from `/specs/002-now-command/spec.md`

## Summary

Add the `ua now` command to the CLI to provide students with an instant, sub-second lookup of their current or immediate next class based on a 30-minute threshold rule. 

## Technical Context

**Language/Version**: Go 1.25+  
**Primary Dependencies**: Cobra (CLI routing), Lipgloss (TUI cards/styling)  
**Storage**: Existing JSON file cache (`~/.ua-cli/cache/schedule.json`)  
**Testing**: Go standard testing (`testing` package)  
**Target Platform**: Windows, macOS, Linux (CLI)  
**Project Type**: Single CLI Application  
**Performance Goals**: < 0.5 seconds execution time  
**Constraints**: Must run entirely offline if the local schedule cache is less than the expiration threshold. Must handle timezones dynamically based on the local system time.  
**Scale/Scope**: O(N) where N is the number of classes in a week (typically < 30). Filtering overhead is negligible.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Speed is a Feature**: Command must run instantly without network overhead if cache is valid. Check passes.
- **UX terminal-first**: Command output will use Lipgloss to render a clean, readable card rather than a full dense table. Check passes.
- **Arquitectura hexagonal**: Logic will be kept in `cmd/ua-cli/now.go` delegating data fetching to the existing `ScheduleService`. Check passes.

## Project Structure

### Documentation (this feature)

```text
specs/002-now-command/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── contracts/           # Internal Contracts
```

### Source Code (repository root)

```text
ua-cli/
├── cmd/
│   └── ua-cli/
│       └── now.go        # Cobra command registration and formatting
├── internal/
│   └── domain/
│       └── schedule/     # (Existing structs used by ScheduleService)
└── tests/
```

**Structure Decision**: Since we are adding a single command that relies entirely on an existing data adapter (the UACloud Schedule JSON parser), we only need to add `now.go` into the main `cmd/ua-cli/` package. The core domain is unchanged.
