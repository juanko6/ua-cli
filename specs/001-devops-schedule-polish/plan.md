# Implementation Plan: DevOps Infrastructure & Schedule Polish

**Branch**: `001-devops-schedule-polish` | **Date**: 2026-02-17 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-devops-schedule-polish/spec.md`

## Summary

Set up GitHub Flow, Semantic Versioning, CI/CD (GitHub Actions + GoReleaser) and fix three schedule bugs: week filtering, event type classification, and deduplication.

## Technical Context

**Language/Version**: Go 1.24 (auto-upgraded via toolchain)  
**Primary Dependencies**: Cobra (CLI), Bubbletea/Lipgloss (TUI), golangci-lint (linting)  
**Storage**: JSON file cache (`~/.ua-cli/cache/schedule.json`)  
**Testing**: `go test` with standard library  
**Target Platform**: Linux amd64, macOS arm64, Windows amd64  
**Project Type**: Single CLI application  
**Performance Goals**: `ua schedule` < 1s (cache hit), < 3s (network fetch)  
**Constraints**: Offline-capable (serve from cache), < 20MB binary  
**Scale/Scope**: Single user, ~100 events/semester

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Speed | ✅ PASS | CI < 5min. schedule filter is O(n), negligible |
| II. UX First | ✅ PASS | TUI preserved. Text mode for pipes. Event types improve readability |
| III. Privacy | ✅ PASS | No telemetry. Cookie stays local. CI doesn't access credentials |
| IV. Maintainability | ✅ PASS | Hexagonal preserved. Filter/dedup logic in domain/service layer |
| V. Resilience | ✅ PASS | Dedup + filter protect against API changes |

**GATE RESULT**: ✅ All gates pass.

## Project Structure

### Documentation (this feature)

```text
specs/001-devops-schedule-polish/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
ua-cli/
├── .github/
│   └── workflows/
│       ├── ci.yml           # [NEW] Lint + Build + Test
│       └── release.yml      # [NEW] GoReleaser on v* tags
├── .goreleaser.yml          # [NEW] Cross-compile config
├── CHANGELOG.md             # [NEW] Keep a Changelog format
├── cmd/ua-cli/
│   ├── main.go
│   ├── root.go
│   └── schedule.go          # [MODIFY] Trim newline from cookie
├── internal/
│   ├── domain/schedule/
│   │   └── entity.go        # [MODIFY] Add Equals() for dedup
│   ├── adapters/
│   │   ├── uacloud/
│   │   │   └── parser.go    # [MODIFY] Improve color→type mapping
│   │   └── presenter/
│   │       ├── text.go
│   │       ├── tui.go
│   │       └── json.go
│   └── service/schedule/
│       └── service.go       # [MODIFY] Add week filter + dedup logic
```

**Structure Decision**: Existing Hexagonal layout preserved. New files only for CI/CD config.

## Complexity Tracking

No constitution violations. No complexity justification needed.
