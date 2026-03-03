# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-03-03

### Added
- New `ua login` command with interactive browser-assisted flow.
- Manual login mode via `ua login --cookie` for headless/SSH usage.
- Session status check via `ua login --status`.
- Auth domain and service layers (`Session`, `CredentialStore`, `AuthService`).
- File credential store with strict `0600` permission enforcement.
- Cookie validation adapter and cross-platform browser opener.
- Root command auth middleware (`PersistentPreRunE`) with auth-exempt command list.
- Basic tests for auth store/service and auth exemption behavior.

### Changed
- `ua` help/description now includes smart-login behavior.
- `ua schedule` now relies on shared auth cookie store instead of manual cookie file handling.

### Fixed
- Edge case: cookie files with incorrect permissions are detected as invalid sessions.

## [0.2.0] - 2026-02-18

### Added
- CI pipeline (GitHub Actions): lint, build, test on push/PR.
- Release pipeline: GoReleaser cross-compiles binaries on `v*` tags.
- `.goreleaser.yml` for Linux/macOS/Windows builds.
- `--next` / `--prev` flags for week navigation in `ua schedule`.

### Fixed
- Schedule now filters events to the current week only (was showing entire semester).
- Event type classification improved with expanded color mapping + className fallback.
- Duplicate events removed using title+start+end composite key deduplification.
- Cookie file reading now trims trailing newline/whitespace.

### Changed
- Default branch renamed from `master` to `main`.

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
