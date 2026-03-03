# Research: Automatic Cookie Login

**Feature**: 001-auto-cookie-login  
**Date**: 2026-03-03

## R1: Cookie Capture Mechanism for CAS SSO

**Decision**: Local HTTP reverse proxy using Go's `net/http/httputil.ReverseProxy`

**Rationale**:
- CAS SSO sets `HttpOnly` cookies — JavaScript extraction is impossible
- No OAuth2/callback mechanism available — CAS uses form-based login and server-side redirects
- A reverse proxy transparently intercepts all `Set-Cookie` headers via `ModifyResponse`
- Zero external dependencies — uses Go standard library only (`net/http`, `net/http/httputil`)

**Alternatives considered**:
| Alternative | Rejected Because |
|-------------|-----------------|
| Browser extension | Requires user to install extension; not cross-platform |
| CDP (Chrome DevTools Protocol) | Couples to Chrome; heavyweight dependency; not portable |
| Direct `document.cookie` | HttpOnly cookies are inaccessible to JavaScript |
| Local callback server (OAuth-style) | CAS doesn't support custom redirect URIs; the `service` param is validated |
| Read browser cookie database | OS/browser-specific; requires decryption; fragile |

---

## R2: CAS Redirect Flow Handling

**Decision**: Rewrite `Location` headers and CAS `service` parameter to route through localhost proxy

**Rationale**:
- CAS SSO flow: UACloud → CAS login → CAS callback → UACloud (with Set-Cookie)
- The proxy must intercept the final redirect where UACloud sets `.ASPXFORMSAUTH*` cookies
- By rewriting `Location` headers from `cvnet.cpd.ua.es` to `localhost:<port>`, the entire flow stays within the proxy
- The CAS `service` URL parameter must also be rewritten so CAS redirects back to the proxy

**Risk**: CAS may validate the `service` URL against a whitelist. If rejected:
- **Mitigation 1**: Proxy only needs to capture cookies from UACloud's `Set-Cookie` headers; CAS itself doesn't set the needed cookies
- **Mitigation 2**: If CAS rejects, fall back to prompting the user for manual cookie paste with clear instructions

---

## R3: Port Selection Strategy

**Decision**: Start with port 18923 (uncommon), fall back to OS-assigned port if busy

**Rationale**:
- Using a fixed high port avoids conflict with common services
- `net.Listen("tcp", ":0")` can auto-assign an available port as fallback
- The URL opened in the browser will use whatever port was actually bound

---

## R4: Success Page & Shutdown

**Decision**: Serve an HTML success page, then gracefully shut down after 2-second delay

**Rationale**:
- The user needs visual confirmation in the browser that login succeeded
- A short delay ensures the response is fully sent before the server closes
- `context.WithCancel` signals the main goroutine to stop blocking and return cookies
