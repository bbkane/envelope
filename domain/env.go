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

// -- EnvVar

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

// -- interface

type EnvService interface {
	CreateEnv(ctx context.Context, args CreateEnvArgs) (*Env, error)
	UpdateEnv(ctx context.Context, name string, args UpdateEnvArgs) error

	CreateLocalEnvVar(ctx context.Context, args CreateLocalEnvVarArgs) (*LocalEnvVar, error)
	ListLocalEnvVars(ctx context.Context, envName string) ([]LocalEnvVar, error)
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
