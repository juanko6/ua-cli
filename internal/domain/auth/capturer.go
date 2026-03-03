package auth

import "context"

// CookieCapturer abstracts browser-based cookie capture.
// Implementations start a local process (e.g. reverse proxy),
// open the browser, and return the captured cookie string.
type CookieCapturer interface {
	// Capture starts the login flow and blocks until cookies are
	// obtained or the context deadline is exceeded.
	Capture(ctx context.Context) (cookie string, err error)
}
