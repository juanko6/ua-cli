package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultValidationURL = "https://cvnet.cpd.ua.es/uaCloud"

func ValidateCookie(cookie string) error {
	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return fmt.Errorf("cookie is empty")
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, defaultValidationURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "ua-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("validate cookie request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusFound {
		return fmt.Errorf("session cookie rejected (status %d)", resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("session validation failed (status %d)", resp.StatusCode)
	}
	return nil
}
