package auth

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFileCredentialStore_SaveLoadAndPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cookie.txt")
	store := NewFileCredentialStore(path)

	cookie := "sessionid=abc123; other=1"
	if err := store.Save(cookie); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded != cookie {
		t.Fatalf("loaded cookie mismatch: %q != %q", loaded, cookie)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if runtime.GOOS != "windows" {
		if got := fi.Mode().Perm(); got != 0o600 {
			t.Fatalf("expected 0600 perms, got %o", got)
		}
	}
}
