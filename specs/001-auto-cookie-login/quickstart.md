# Quickstart: Automatic Cookie Login

**Feature**: 001-auto-cookie-login  
**Date**: 2026-03-03

## Prerequisites

- Go 1.25+ installed
- Access to UA university credentials

## Build

```bash
go build -o ua.exe ./cmd/ua-cli/
```

## Usage

### Automatic login (default)

```bash
ua login
```

This will:
1. Start a local proxy server
2. Open your browser to the UACloud login page
3. After you log in, cookies are captured automatically
4. CLI displays "✓ Login successful"

### Manual fallback

```bash
ua login --cookie
```

Paste your cookie string from browser DevTools.

### Check session status

```bash
ua login --status
```

## How It Works

```
ua login → starts localhost proxy → opens browser
       ↓
User logs in via CAS SSO in browser
       ↓
Proxy intercepts Set-Cookie headers from UACloud
       ↓
Cookies saved to ~/.ua-cli/cookie.txt
       ↓
CLI reports success, proxy shuts down
```

## Verification

```bash
# Test that login saved a valid session
ua login --status
# Expected: ✓ Session active

# Test that schedule works with the new session
ua schedule
```
