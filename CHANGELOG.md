# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-17

### Added
- Initial CLI skeleton with Cobra.
- `ua schedule` command with adaptive TUI/text/JSON output.
- UACloud schedule API integration with cookie authentication.
- JSON file cache for offline-first schedule access.
- Bubbletea interactive TUI for terminal mode.
- Plain text table output for piped commands.
- JSON output format via `--json` flag.
- Hexagonal architecture (domain, service, adapters, presenters).
