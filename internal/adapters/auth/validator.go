package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const defaultValidationURL = "https://cvnet.cpd.ua.es/uaCloud"

func ValidateCookie(cookie string) error {
	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return fmt.Errorf("cookie is empty")
	}

	// Build a cookie jar so cookies are forwarded through CAS redirects.
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://cvnet.cpd.ua.es")
	var cookies []*http.Cookie
	for _, raw := range strings.Split(cookie, ";") {
		parts := strings.SplitN(strings.TrimSpace(raw), "=", 2)
		if len(parts) == 2 {
			cookies = append(cookies, &http.Cookie{Name: parts[0], Value: parts[1]})
		}
	}
	jar.SetCookies(u, cookies)

	client := &http.Client{
		Jar:     jar,
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, defaultValidationURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "ua-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("validate cookie request: %w", err)
	}
	defer resp.Body.Close()

	// After following all redirects, check where we ended up.
	// Valid session → final URL is on cvnet.cpd.ua.es
	// Invalid session → final URL stays on autentica.cpd.ua.es (CAS login page)
	finalURL := resp.Request.URL.String()
	if strings.Contains(finalURL, "autentica.cpd.ua.es") {
		return fmt.Errorf("session cookie rejected: redirected to CAS login")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("session cookie rejected (status %d)", resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("session validation failed (status %d)", resp.StatusCode)
	}
	return nil
}
