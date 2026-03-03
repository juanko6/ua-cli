package auth

// CredentialStore abstracts cookie persistence.
type CredentialStore interface {
	Load() (string, error)
	Save(cookie string) error
	Exists() (bool, error)
	Delete() error
}
