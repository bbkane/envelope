package domain

import (
	"context"
	"time"
)

// -- Env

type Env struct {
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
}

type EnvCreateArgs struct {
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
}

type EnvUpdateArgs struct {
	Comment    *string
	CreateTime *time.Time
	NewName    *string
	UpdateTime *time.Time
}

// -- EnvLocalVar

type EnvLocalVar struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type EnvLocalVarCreateArgs struct {
	EnvName    string
	Name       string
	Comment    string
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

// -- interface

type EnvService interface {
	EnvCreate(ctx context.Context, args EnvCreateArgs) (*Env, error)
	EnvDelete(ctx context.Context, name string) error
	EnvUpdate(ctx context.Context, name string, args EnvUpdateArgs) error
	EnvShow(ctx context.Context, name string) (*Env, error)

	EnvLocalVarCreate(ctx context.Context, args EnvLocalVarCreateArgs) (*EnvLocalVar, error)
	EnvLocalVarDelete(ctx context.Context, envName string, name string) error
	EnvLocalVarList(ctx context.Context, envName string) ([]EnvLocalVar, error)
	EnvLocalVarShow(ctx context.Context, envName string, name string) (*EnvLocalVar, error)

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

// StringToTimeMust works like StringToTime but panics on errors.
// I think this is usually acceptable as times are formatted pretty carefully
// in the db
func StringToTimeMust(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.Round(0)
}
