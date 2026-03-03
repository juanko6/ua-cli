# Data Model: Automatic Cookie Login

**Feature**: 001-auto-cookie-login  
**Date**: 2026-03-03

## Entities

### LoginProxy

Temporary local reverse proxy that intercepts cookies during browser-based CAS SSO login.

| Field | Type | Description |
|-------|------|-------------|
| Port | integer | TCP port the proxy listens on |
| TargetHost | string | UACloud base URL (`https://cvnet.cpd.ua.es`) |
| Timeout | duration | Maximum wait time before auto-cancellation (2 min) |
| CapturedCookies | string | The raw cookie string captured from Set-Cookie headers |
| Status | enum | `waiting`, `captured`, `timeout`, `error` |

**Validation Rules**:
- Port must be between 1024 and 65535
- Timeout must be positive and ≤ 5 minutes
- CapturedCookies must contain at least `.ASPXFORMSAUTH` to be considered valid

**State Transitions**:
```
waiting → captured    (Set-Cookie with auth cookies detected)
waiting → timeout     (2-minute deadline reached)
waiting → error       (port bind failure, network error)
```

---

### Session (existing entity — unchanged)

| Field | Type | Description |
|-------|------|-------------|
| Status | enum | `missing`, `valid`, `expired`, `invalid` |
| Cookie | string | The stored cookie string |
| Age | duration | Time since cookie was saved |
| Path | string | Path to the cookie file on disk |

No changes to the existing Session entity. The LoginProxy produces a cookie string that is stored through the existing `CredentialStore` interface.
