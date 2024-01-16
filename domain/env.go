package domain

import (
	"context"
	"time"
)

// -- Env

type Env struct {
	Name       string
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
}

type CreateEnvArgs struct {
	Name       string
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
}

type UpdateEnvArgs struct {
	Comment    *string
	CreateTime *time.Time
	NewName    *string
	UpdateTime *time.Time
}

// -- LocalEnvVar

type LocalEnvVar struct {
	EnvName    string
	Name       string
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type CreateLocalEnvVarArgs struct {
	EnvName    string
	Name       string
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

// -- Keyring

// Keyring provides a simple set/get interface for a keyring service.
type Keyring interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Service() string
}

type KeyringEntry struct {
	Name       string
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type KeyringEntryCreateArgs struct {
	Name       string
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

// -- interface

type EnvService interface {
	EnvCreate(ctx context.Context, args CreateEnvArgs) (*Env, error)
	EnvUpdate(ctx context.Context, name string, args UpdateEnvArgs) error

	EnvVarLocalCreate(ctx context.Context, args CreateLocalEnvVarArgs) (*LocalEnvVar, error)
	EnvVarLocalList(ctx context.Context, envName string) ([]LocalEnvVar, error)

	KeyringEntryCreate(ctx context.Context, args KeyringEntryCreateArgs) (*KeyringEntry, error)
}

// -- Utility function

// TimeToString converts a time to UTC, then formats as RFC3339
func TimeToString(t time.Time) string {

	return t.Round(0).UTC().Format(time.RFC3339)
}

// StringToTime converts a RFC3339 formatted string into a time.Time
func StringToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return t, err
	}
	return t.Round(0), nil
}
