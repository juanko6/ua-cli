package auth

import "time"

// SessionStatus describes the state of the saved UA session cookie.
type SessionStatus string

const (
	SessionValid   SessionStatus = "valid"
	SessionExpired SessionStatus = "expired"
	SessionMissing SessionStatus = "missing"
	SessionInvalid SessionStatus = "invalid"
)

// Session represents a persisted login session.
type Session struct {
	Status    SessionStatus
	Cookie    string
	Path      string
	UpdatedAt time.Time
	Age       time.Duration
}
