# UA-CLI Constitution

<!-- HTML Comment: Sync Impact Report -->
<!-- Version: 1.0.0 (Initial Ratification) -->
<!-- Added Principles: Speed, UX, Privacy, Maintainability, Resilience -->

## Core Principles

### I. Speed is a Feature
Commands **MUST** execute in milliseconds whenever possible. We prioritize aggressive caching strategy (Offline First) to ensure instant feedback, masking network latency. The user should never wait for a command that can be served from local cache.

### II. User Experience (UX) First
We embrace "Unhinged Aesthetics" using modern TUI libraries (Bubbletea/Lipgloss). Interactions **MUST** be delightful, interactive, and visually polished. However, standard stdout/stderr predictability **MUST** be preserved for scripting compatibility (pipe-friendly when not in TTY).

### III. Privacy & Security
Data is **Local-First**. Credentials **MUST** be stored securely using the OS native Keyring or encrypted cookies. We **MUST NOT** implement any third-party tracking or telemetry. User data never leaves their machine except to communicate with official UA servers.

### IV. Maintainability (Hexagonal)
We strictly follow **Hexagonal Architecture (Ports & Adapters)** and the **Standard Go Project Layout**. Business logic (Domain) **MUST** be decoupled from infrastructure (HTTP, CLI, TUI). High test coverage is mandatory, especially for HTML parsers.

### V. Resilience
Scraping logic **MUST** be decoupled and resilient. Parsers should be robust against minor HTML structure changes. We acknowledge the risk of "Institutional Blocking" and mitigate it with respectful rate limiting and configurable User-Agents.

## Feature Scope
**Included**: Authentication (Keyring/Cookie), Schedule, Grades, Notices, Virtual Campus Download.
**Excluded**: Full Offline Mode (write), Push Notifications (system), Advanced File Organization.

## Governance

This Constitution supersedes all other technical decisions.
Amendments require a Pull Request with a clear rationale and must not violate the core values of Privacy or Speed without overwhelming justification.

**Version**: 1.0.0 | **Ratified**: 2026-02-17 | **Last Amended**: 2026-02-17
