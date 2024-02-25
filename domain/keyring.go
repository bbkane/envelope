package domain

import (
	"errors"
	"time"
)

// -- Keyring

var (
	// ErrKeyringSecretNotFound is the expected error if the secret isn't found in the
	// keyring.
	ErrKeyringSecretNotFound = errors.New("secret not found in keyring")
	// ErrKeyringSetDataTooBig is returned if `Set` was called with too much data.
	// On MacOS: The combination of service, username & password should not exceed ~3000 bytes
	// On Windows: The service is limited to 32KiB while the password is limited to 2560 bytes
	// On Linux/Unix: There is no theoretical limit but performance suffers with big values (>100KiB)
	ErrKeyringSetDataTooBig = errors.New("data passed to Set was too big")
)

// Keyring provides a simple set/get interface for a keyring service.
type Keyring interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Service() string
}

type KeyringEntry struct {
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type KeyringEntryCreateArgs struct {
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}
