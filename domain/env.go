package domain

import (
	"context"
	"time"
)

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

type EnvService interface {
	CreateEnv(ctx context.Context, args CreateEnvArgs) (*Env, error)
	UpdateEnv(ctx context.Context, name string, args UpdateEnvArgs) error
}

// TimeToString converts a time to UTC, then formats as RFC3339
func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// StringToTime converts a RFC3339 formatted string into a time.Time
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
