package domain

import (
	"context"
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

type EnvService interface {
	CreateEnv(ctx context.Context, args CreateEnvArgs) (EnvID, error)
}
