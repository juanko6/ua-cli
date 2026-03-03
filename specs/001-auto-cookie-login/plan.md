# Implementation Plan: Automatic Cookie Login

**Branch**: `001-auto-cookie-login` | **Date**: 2026-03-03 | **Spec**: [spec.md](file:///c:/Users/juankodev/Desktop/ua-cli/specs/001-auto-cookie-login/spec.md)  
**Input**: Feature specification from `/specs/001-auto-cookie-login/spec.md`

## Summary

Replace the manual cookie copy-paste login flow with an automatic browser-based login. The CLI starts a local reverse proxy (`httputil.ReverseProxy`) that proxies `cvnet.cpd.ua.es`, opens the user's browser through it, and intercepts `Set-Cookie` headers from the CAS SSO flow using `ModifyResponse`. The captured cookies are stored via the existing `FileCredentialStore`. The `--cookie` flag remains as manual fallback.

## Technical Context

**Language/Version**: Go 1.25.6  
**Primary Dependencies**: Standard library (`net/http`, `net/http/httputil`, `crypto/tls`), Cobra, Lipgloss (existing)  
**Storage**: File-based cookie storage at `~/.ua-cli/cookie.txt` (existing mechanism, unchanged)  
**Testing**: `go test ./...`  
**Target Platform**: Windows, macOS, Linux (cross-platform CLI)  
**Project Type**: Single CLI project  
**Performance Goals**: Proxy startup < 500ms; cookie capture immediate upon CAS redirect  
**Constraints**: No external dependencies beyond existing go.mod; 2-minute login timeout  
**Scale/Scope**: Single-user CLI tool; one concurrent login session at a time

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Speed is a Feature | ✅ PASS | Proxy starts in <500ms; no impact on cached commands |
| II. UX First | ✅ PASS | Eliminates manual DevTools copy-paste; uses Lipgloss for output styling |
| III. Privacy & Security | ✅ PASS | Cookies stored locally via existing FileCredentialStore; proxy runs on localhost only; no data sent to third parties |
| IV. Maintainability (Hexagonal) | ✅ PASS | New `CookieCapturer` port/interface; `ProxyAdapter` implements it; cleanly separated from login command logic |
| V. Resilience | ✅ PASS | Graceful timeout; port fallback; manual `--cookie` fallback if proxy fails |
| Feature Scope | ✅ PASS | Authentication is explicitly in scope |

**All gates pass. No violations.**

## Project Structure

### Documentation (this feature)

```text
specs/001-auto-cookie-login/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0: research decisions
├── data-model.md        # Phase 1: entity definitions
├── quickstart.md        # Phase 1: usage guide
├── contracts/           # Phase 1: internal port contracts
│   └── internal-ports.md
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/
└── ua-cli/
    └── login.go              # MODIFY: wire proxy-based flow as default

internal/
├── adapters/
│   └── auth/
│       ├── browser.go        # EXISTING: opens browser (reused)
│       ├── proxy.go          # NEW: reverse proxy cookie capturer
│       ├── validator.go      # EXISTING: cookie validation (already fixed)
│       ├── filestore.go      # EXISTING: cookie persistence (unchanged)
│       └── proxy_test.go     # NEW: proxy unit tests
├── domain/
│   └── auth/
│       ├── capturer.go       # NEW: CookieCapturer port interface
│       ├── session.go        # EXISTING: unchanged
│       └── store.go          # EXISTING: unchanged
└── service/
    └── auth/
        └── service.go        # EXISTING: minor changes to integrate capturer
```

**Structure Decision**: Follows existing Hexagonal Architecture layout. New port (`CookieCapturer`) in `domain/auth/`, new adapter (`ProxyAdapter`) in `adapters/auth/`. CLI wiring in `cmd/ua-cli/login.go`.
