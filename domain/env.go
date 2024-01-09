package domain

import (
	"context"
	"errors"
	"time"
)

type EnvID int64

type Env struct {
	ID         EnvID
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
	Name       *string
	Comment    *string
	CreateTime *time.Time
	UpdateTime *time.Time
}

type EnvService interface {
	CreateEnv(ctx context.Context, args CreateEnvArgs) (EnvID, error)
	UpdateEnv(ctx context.Context, args UpdateEnvArgs) error
}

// TimeToString converts a time to UTC, then formats as RFC3339
func TimeToString(t time.Time) (string, error) {
	if t.IsZero() {
		return "", errors.New("time should not be 0")
	}
	return t.UTC().Format(time.RFC3339), nil
}

// StringToTime converts a RFC3339 formatted string into a time.Time
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
