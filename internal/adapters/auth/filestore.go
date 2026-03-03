package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DefaultCookieRelativePath = ".ua-cli/cookie.txt"

type FileCredentialStore struct {
	path string
}

func NewFileCredentialStore(path string) *FileCredentialStore {
	return &FileCredentialStore{path: path}
}

func DefaultCookiePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, DefaultCookieRelativePath), nil
}

func (s *FileCredentialStore) Path() string {
	return s.path
}

func (s *FileCredentialStore) Load() (string, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func (s *FileCredentialStore) Save(cookie string) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(s.path, []byte(strings.TrimSpace(cookie)), 0o600); err != nil {
		return err
	}
	if err := os.Chmod(s.path, 0o600); err != nil {
		return err
	}
	return nil
}

func (s *FileCredentialStore) Exists() (bool, error) {
	_, err := os.Stat(s.path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *FileCredentialStore) Delete() error {
	err := os.Remove(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *FileCredentialStore) Stat() (os.FileInfo, error) {
	fi, err := os.Stat(s.path)
	if err != nil {
		return nil, fmt.Errorf("stat cookie file: %w", err)
	}
	return fi, nil
}
