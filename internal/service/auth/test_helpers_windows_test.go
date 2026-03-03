//go:build windows

package auth

import (
	"os"
	"time"
)

func osChtimes(path string, when time.Time) error {
	return os.Chtimes(path, when, when)
}
