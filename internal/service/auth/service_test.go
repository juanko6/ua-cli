package auth

import (
	"testing"
	"time"

	adaptauth "ua-cli/internal/adapters/auth"
	domainauth "ua-cli/internal/domain/auth"
)

func TestCheckSessionMissing(t *testing.T) {
	path := t.TempDir() + "/cookie.txt"
	store := adaptauth.NewFileCredentialStore(path)
	svc := NewAuthService(store, nil)

	s, err := svc.CheckSession()
	if err != nil {
		t.Fatalf("check session: %v", err)
	}
	if s.Status != domainauth.SessionMissing {
		t.Fatalf("expected missing, got %s", s.Status)
	}
}

func TestCheckSessionExpiredByAge(t *testing.T) {
	path := t.TempDir() + "/cookie.txt"
	store := adaptauth.NewFileCredentialStore(path)
	if err := store.Save("session=ok"); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-(ExpiredAgeThreshold + time.Hour))
	if err := osChtimes(path, old); err != nil {
		t.Fatal(err)
	}

	svc := NewAuthService(store, func(cookie string) error { return nil })
	s, err := svc.CheckSession()
	if err != nil {
		t.Fatalf("check session: %v", err)
	}
	if s.Status != domainauth.SessionExpired {
		t.Fatalf("expected expired, got %s", s.Status)
	}
}
