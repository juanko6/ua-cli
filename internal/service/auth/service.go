package auth

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	adaptauth "ua-cli/internal/adapters/auth"
	domainauth "ua-cli/internal/domain/auth"
)

const (
	WarningAgeThreshold = time.Hour
	ExpiredAgeThreshold = 24 * time.Hour
)

type CookieValidator func(cookie string) error

type AuthService struct {
	store     domainauth.CredentialStore
	validator CookieValidator
}

func NewAuthService(store domainauth.CredentialStore, validator CookieValidator) *AuthService {
	return &AuthService{store: store, validator: validator}
}

func (s *AuthService) CheckSession() (*domainauth.Session, error) {
	exists, err := s.store.Exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return &domainauth.Session{Status: domainauth.SessionMissing}, nil
	}

	cookie, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(cookie) == "" {
		return &domainauth.Session{Status: domainauth.SessionInvalid}, nil
	}

	session := &domainauth.Session{Status: domainauth.SessionValid, Cookie: cookie}

	if fs, ok := s.store.(*adaptauth.FileCredentialStore); ok {
		fi, err := fs.Stat()
		if err == nil {
			session.Path = fs.Path()
			session.UpdatedAt = fi.ModTime()
			session.Age = time.Since(fi.ModTime())
			if runtime.GOOS != "windows" && fi.Mode().Perm() != 0o600 {
				session.Status = domainauth.SessionInvalid
				return session, fmt.Errorf("cookie file permissions are %o (expected 600)", fi.Mode().Perm())
			}
		}
	}

	if session.Age > ExpiredAgeThreshold {
		session.Status = domainauth.SessionExpired
		return session, nil
	}

	if s.validator != nil {
		if err := s.validator(cookie); err != nil {
			session.Status = domainauth.SessionExpired
			return session, nil
		}
	}

	return session, nil
}

func (s *AuthService) StoreCookie(cookie string) error {
	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return fmt.Errorf("cookie cannot be empty")
	}
	if s.validator != nil {
		if err := s.validator(cookie); err != nil {
			return err
		}
	}
	return s.store.Save(cookie)
}

func (s *AuthService) LoadCookie() (string, error) {
	return s.store.Load()
}

func (s *AuthService) WarnIfAged(session *domainauth.Session) {
	if session == nil {
		return
	}
	if session.Age > WarningAgeThreshold && session.Status == domainauth.SessionValid {
		fmt.Fprintf(os.Stderr, "Warning: your UA session cookie is %s old and may expire soon.\n", session.Age.Round(time.Minute))
	}
}
