package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

const (
	targetHost    = "cvnet.cpd.ua.es"
	casHost       = "autentica.cpd.ua.es"
	targetScheme  = "https"
	defaultPort   = 18923
	proxyTimeout  = 2 * time.Minute
	shutdownDelay = 2 * time.Second
)

// authCookiePrefixes are the cookie name prefixes we need to capture.
var authCookiePrefixes = []string{
	".ASPXFORMSAUTH",
	"ASP.NET_SessionId",
}

// subAppPaths are UACloud sub-app URLs to visit after initial CAS login
// so the proxy captures each sub-app's .ASPXFORMSAUTH cookie.
var subAppPaths = []string{
	"/uaHorarios",
}

// ProxyAdapter implements domain/auth.CookieCapturer using a local
// reverse proxy that intercepts Set-Cookie headers from UACloud/CAS.
type ProxyAdapter struct {
	Port int // Bound port (set after Start)
}

// Capture starts a local reverse proxy, waits for authentication cookies
// from UACloud, and returns them as a raw cookie header string.
// It blocks until cookies are captured or the context is cancelled.
func (p *ProxyAdapter) Capture(ctx context.Context) (string, error) {
	// Pick a port: try default first, then OS-assigned.
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			return "", fmt.Errorf("listen: %w", err)
		}
	}
	p.Port = ln.Addr().(*net.TCPAddr).Port

	// Accumulate cookies across multiple responses.
	var mu sync.Mutex
	collected := make(map[string]string) // name → "name=value"
	cookieCh := make(chan string, 1)
	done := false
	subAppIdx := 0 // next sub-app to visit

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Route by path: /cas/* → CAS server, everything else → UACloud.
			if strings.HasPrefix(req.URL.Path, "/cas") {
				req.URL.Scheme = targetScheme
				req.URL.Host = casHost
				req.Host = casHost
			} else {
				req.URL.Scheme = targetScheme
				req.URL.Host = targetHost
				req.Host = targetHost

				// Also capture auth cookies from the browser's Cookie header.
				// ASP.NET_SessionId is set once and sent on all requests,
				// but never re-set via Set-Cookie on subsequent responses.
				mu.Lock()
				if !done {
					for _, c := range req.Cookies() {
						for _, prefix := range authCookiePrefixes {
							if strings.HasPrefix(c.Name, prefix) {
								collected[c.Name] = c.Name + "=" + c.Value
								break
							}
						}
					}
				}
				mu.Unlock()
			}
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		},
		ModifyResponse: func(resp *http.Response) error {
			// Strip Domain/Secure/SameSite from Set-Cookie so cookies
			// work on localhost (non-HTTPS).
			cleaned := make([]string, 0, len(resp.Header["Set-Cookie"]))
			for _, sc := range resp.Header["Set-Cookie"] {
				sc = stripCookieAttr(sc, "domain")
				sc = stripCookieAttr(sc, "secure")
				sc = strings.Replace(sc, "SameSite=None", "SameSite=Lax", 1)
				cleaned = append(cleaned, sc)
			}
			resp.Header["Set-Cookie"] = cleaned

			// Rewrite Location headers — replace BOTH CAS and cvnet
			// domains with localhost so the browser stays on the proxy.
			// The service= param stays URL-encoded pointing at cvnet
			// (CAS validates it), only the bare host is rewritten.
			if loc := resp.Header.Get("Location"); loc != "" {
				newLoc := rewriteURL(loc, p.Port)
				// Only strip gateway=true before the initial login.
				// After that, keep it so CAS does a silent pass-through
				// for sub-app cookie capture.
				if !hasASPXFormsAuth(collected) {
					newLoc = stripGateway(newLoc)
				}
				resp.Header.Set("Location", newLoc)
			}

			// Accumulate auth cookies.
			mu.Lock()
			defer mu.Unlock()
			if done {
				return nil
			}

			for _, sc := range resp.Header["Set-Cookie"] {
				nameVal := strings.SplitN(sc, ";", 2)[0]
				name := strings.SplitN(nameVal, "=", 2)[0]
				for _, prefix := range authCookiePrefixes {
					if strings.HasPrefix(name, prefix) {
						collected[name] = nameVal
						break
					}
				}
			}

			// Check if we have .ASPXFORMSAUTH cookies.
			if hasASPXFormsAuth(collected) {
				// If there are sub-apps still to visit, redirect there
				// to capture their cookies too.
				if subAppIdx < len(subAppPaths) {
					nextPath := subAppPaths[subAppIdx]
					subAppIdx++
					resp.Header.Set("Location",
						fmt.Sprintf("http://localhost:%d%s", p.Port, nextPath))
					resp.StatusCode = http.StatusFound
					resp.Status = "302 Found"
				} else if isLoginComplete(collected) {
					// All required cookies captured — signal completion.
					done = true
					var parts []string
					for _, v := range collected {
						parts = append(parts, v)
					}
					select {
					case cookieCh <- strings.Join(parts, "; "):
					default:
					}
					resp.Header.Set("Location",
						fmt.Sprintf("http://localhost:%d/_ua_login_done", p.Port))
					resp.StatusCode = http.StatusFound
					resp.Status = "302 Found"
				}
			}
			return nil
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	mux.HandleFunc("/_ua_login_done", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, successPage)
	})

	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()

	// Wait for cookies or timeout.
	var cookie string
	select {
	case cookie = <-cookieCh:
	case <-ctx.Done():
		_ = srv.Close()
		return "", fmt.Errorf("login timed out — no cookies captured within the deadline")
	}

	// Give the browser time to render the success page, then shut down.
	time.AfterFunc(shutdownDelay, func() { _ = srv.Close() })
	return cookie, nil
}

// ProxyURL returns the URL the browser should be opened to.
func (p *ProxyAdapter) ProxyURL() string {
	return fmt.Sprintf("http://localhost:%d/uaCloud", p.Port)
}

// rewriteURL replaces both cvnet AND CAS host with localhost,
// but only the bare host — URL-encoded instances in query params are kept.
func rewriteURL(rawURL string, port int) string {
	lh := fmt.Sprintf("localhost:%d", port)
	u := strings.Replace(rawURL, "https://"+targetHost, "http://"+lh, 1)
	u = strings.Replace(u, "http://"+targetHost, "http://"+lh, 1)
	u = strings.Replace(u, "https://"+casHost, "http://"+lh, 1)
	u = strings.Replace(u, "http://"+casHost, "http://"+lh, 1)
	return u
}

// stripGateway removes gateway=true from a URL so CAS shows the
// login form instead of silently bouncing back without authentication.
func stripGateway(rawURL string) string {
	// Handle both &gateway=true and ?gateway=true&...
	u := strings.Replace(rawURL, "&gateway=true", "", 1)
	u = strings.Replace(u, "?gateway=true&", "?", 1)
	u = strings.Replace(u, "?gateway=true", "", 1)
	return u
}

// stripCookieAttr removes a Set-Cookie attribute (e.g. "domain", "secure").
func stripCookieAttr(cookie, attr string) string {
	parts := strings.Split(cookie, ";")
	var kept []string
	for _, p := range parts {
		if !strings.HasPrefix(strings.TrimSpace(strings.ToLower(p)), strings.ToLower(attr)) {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, ";")
}

// extractAuthCookies scans Set-Cookie headers for authentication cookies
// and returns them as a semicolon-separated string (format: "key=val; key2=val2").
func extractAuthCookies(setCookies []string) string {
	var parts []string
	for _, sc := range setCookies {
		nameVal := strings.SplitN(sc, ";", 2)[0]
		name := strings.SplitN(nameVal, "=", 2)[0]
		for _, prefix := range authCookiePrefixes {
			if strings.HasPrefix(name, prefix) {
				parts = append(parts, nameVal)
				break
			}
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "; ")
}

// hasASPXFormsAuth checks if at least one .ASPXFORMSAUTH cookie has been collected.
func hasASPXFormsAuth(collected map[string]string) bool {
	for name := range collected {
		if strings.HasPrefix(name, ".ASPXFORMSAUTH") {
			return true
		}
	}
	return false
}

// isLoginComplete checks if all required cookies have been captured,
// including sub-app specific cookies like UAHORARIOS.
func isLoginComplete(collected map[string]string) bool {
	_, hasHorarios := collected[".ASPXFORMSAUTHUAHORARIOS"]
	return hasASPXFormsAuth(collected) && hasHorarios
}
