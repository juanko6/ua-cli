# Internal Contracts: Automatic Cookie Login

This feature is a CLI tool — no REST/GraphQL API contracts apply.
The contracts below define the **internal Go interfaces** between components.

## LoginProxy Port (Domain → Adapter)

```go
// CookieCapturer abstracts browser-based cookie capture.
type CookieCapturer interface {
    // Capture starts the login flow and blocks until cookies are obtained or timeout.
    // Returns the raw cookie string on success.
    Capture(ctx context.Context) (cookie string, err error)
}
```

## CredentialStore Port (existing — unchanged)

```go
// CredentialStore abstracts cookie persistence.
type CredentialStore interface {
    Save(cookie string) error
    Load() (string, os.FileInfo, error)
}
```

## CookieValidator Port (existing — unchanged)

```go
// CookieValidator tests whether a cookie string represents a live session.
type CookieValidator func(cookie string) error
```

## Interaction Flow

```
loginCmd (CLI layer)
  ├─ --cookie flag? → promptCookie() → CredentialStore.Save()
  └─ default       → CookieCapturer.Capture() → CredentialStore.Save()
                         ↓
                    ProxyAdapter (starts localhost proxy)
                         ↓
                    ModifyResponse intercepts Set-Cookie
                         ↓
                    Returns cookie string via channel
```
