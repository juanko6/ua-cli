# Feature Specification: Automatic Cookie Login

**Feature Branch**: `001-auto-cookie-login`  
**Created**: 2026-03-03  
**Status**: Draft  
**Input**: User description: "Automatic cookie extraction for ua login: the CLI should open the browser, let the user log in to UACloud via CAS SSO, and automatically capture the authentication cookies without manual copy-paste"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - One-Command Automatic Login (Priority: P1)

As a university student using the CLI, I want to run `ua login` and have the system automatically open my browser, let me log in to UACloud, and capture my session cookies — so I don't have to manually copy cookies from DevTools.

**Why this priority**: This is the core value proposition. The current manual cookie copy-paste is error-prone (HttpOnly cookies, wrong format, expired sessions) and creates a terrible first-use experience.

**Independent Test**: Can be fully tested by running `ua login`, completing the university login flow in the browser, and verifying the CLI reports success and stores a valid session.

**Acceptance Scenarios**:

1. **Given** the user has no saved session, **When** they run `ua login`, **Then** the system opens their default browser to the UACloud login page, and after completing CAS authentication, the CLI automatically captures the session cookies and displays "✓ Login successful".
2. **Given** the user completes login in the browser, **When** the session cookies are captured, **Then** the cookies are securely stored and usable by all other CLI commands (e.g., `ua schedule`).
3. **Given** the browser is opened for login, **When** the user completes authentication, **Then** the browser displays a success page telling the user they can close the tab.

---

### User Story 2 - Manual Cookie Fallback (Priority: P2)

As a user in a restricted environment (e.g., SSH session, headless server), I want to manually paste my cookie using the `--cookie` flag so I can still authenticate when the automatic browser flow isn't available.

**Why this priority**: Provides a necessary fallback for edge cases where automatic browser capture isn't possible.

**Independent Test**: Can be tested by running `ua login --cookie`, pasting a valid cookie string, and verifying successful session storage.

**Acceptance Scenarios**:

1. **Given** the user runs `ua login --cookie`, **When** they paste a valid cookie string, **Then** the system validates and stores it, displaying "✓ Login successful".

---

### User Story 3 - Login Timeout Protection (Priority: P3)

As a user, I want the login process to automatically cancel after a reasonable timeout so the CLI doesn't hang indefinitely if I close the browser or forget to complete login.

**Why this priority**: Prevents poor UX from abandoned login attempts.

**Independent Test**: Can be tested by running `ua login` and waiting without completing the browser login, verifying the CLI exits with a timeout message.

**Acceptance Scenarios**:

1. **Given** the user runs `ua login`, **When** they do not complete browser authentication within 2 minutes, **Then** the system displays a timeout error and exits cleanly.

---

### Edge Cases

- What happens when the user's browser cannot open automatically? → Fall back to displaying the URL and instructions.
- What happens when the user is already logged in with a valid session? → Inform them the session is already active, offer to re-authenticate.
- What happens when CAS login fails (wrong credentials)? → The system waits for a successful login or times out.
- What happens when the local port needed for cookie capture is already in use? → Try alternative ports automatically.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST start a local process to intercept authentication cookies from the UACloud/CAS SSO login flow.
- **FR-002**: System MUST automatically open the user's default browser to the UACloud login page when `ua login` is executed.
- **FR-003**: System MUST capture the authentication cookies (`.ASPXFORMSAUTH*`, `ASP.NET_SessionId`) after successful CAS login without user intervention.
- **FR-004**: System MUST securely store captured cookies in the same format and location used by the existing cookie storage mechanism.
- **FR-005**: System MUST display a user-friendly success page in the browser after cookies are captured.
- **FR-006**: System MUST terminate the local process and exit cleanly after cookies are captured or after a 2-minute timeout.
- **FR-007**: System MUST preserve the `--cookie` flag for manual cookie input as a fallback method.
- **FR-008**: System MUST display clear progress messages in the terminal during the login flow.
- **FR-009**: System MUST handle port conflicts by trying alternative ports if the default one is occupied.

### Key Entities

- **Session Cookie**: The HTTP cookie string containing authentication tokens for UACloud. Composed of multiple cookie key-value pairs separated by semicolons.
- **Login Proxy**: A temporary local process that intercepts cookies during the CAS authentication flow. Lives only for the duration of the login attempt.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can complete the login process in under 30 seconds (excluding time spent entering university credentials in the browser).
- **SC-002**: 100% of successful CAS logins result in automatically captured and stored cookies, requiring zero manual copy-paste.
- **SC-003**: The login flow works on all supported platforms (Windows, macOS, Linux) without additional user configuration.
- **SC-004**: The fallback manual mode (`--cookie`) continues to function identically to its current behavior.
