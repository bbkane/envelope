package domain

import (
	"context"
	"time"
)

type Env struct {
	ID int64

	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
}

type CreateEnvArgs struct {
	Comment    *string
	CreateTime time.Time
	UpdateTime time.Time
}

type EnvService interface {
	CreateEnv(ctx context.Context, args CreateEnvArgs) (*Env, error)
}
