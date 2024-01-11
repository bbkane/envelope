package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/connect"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

func DerefOrEmpty[T any](val *T) T {
	if val == nil {
		var empty T
		return empty
	}
	return *val
}

func IsNotNil[T any](val *T) bool {
	return val != nil
}

type EnvService struct {
	db *sql.DB
}

func NullStringToStrPtr(val sql.NullString) *string {
	if !val.Valid {
		return nil
	}
	return &val.String
}

func NewEnvService(ctx context.Context, dsn string) (domain.EnvService, error) {
	// TODO use context!!
	db, err := connect.Connect(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}
	return &EnvService{
		db: db,
	}, nil
}

func (e *EnvService) CreateEnv(ctx context.Context, args domain.CreateEnvArgs) (*domain.Env, error) {
	queries := sqlcgen.New(e.db)

	createdEnvID, err := queries.CreateEnv(ctx, sqlcgen.CreateEnvParams{
		Name: args.Name,
		Comment: sql.NullString{
			String: DerefOrEmpty(args.Comment),
			Valid:  IsNotNil(args.Comment),
		},
		CreateTime: domain.TimeToString(args.CreateTime),
		UpdateTime: domain.TimeToString(args.UpdateTime),
	})

	if err != nil {
		return nil, fmt.Errorf("could not create env in db: %w", err)
	}

	createTime, err := domain.StringToTime(createdEnvID.CreateTime)
	if err != nil {
		panic(err)
	}
	updateTime, err := domain.StringToTime(createdEnvID.UpdateTime)
	if err != nil {
		panic(err)
	}

	return &domain.Env{
		Name:       createdEnvID.Name,
		Comment:    NullStringToStrPtr(createdEnvID.Comment),
		CreateTime: createTime,
		UpdateTime: updateTime,
	}, nil
}

func (e *EnvService) UpdateEnv(ctx context.Context, name string, args domain.UpdateEnvArgs) error {

	queries := sqlcgen.New(e.db)

	err := queries.UpdateEnv(ctx, sqlcgen.UpdateEnvParams{
		NewName: sql.NullString{
			String: DerefOrEmpty(args.NewName),
			Valid:  IsNotNil(args.NewName),
		},
		Comment: sql.NullString{
			String: DerefOrEmpty(args.Comment),
			Valid:  IsNotNil(args.Comment),
		},
		CreateTime: sql.NullString{
			String: domain.TimeToString(DerefOrEmpty(args.CreateTime)),
			Valid:  IsNotNil(args.CreateTime),
		},
		UpdateTime: sql.NullString{
			String: domain.TimeToString(DerefOrEmpty(args.CreateTime)),
			Valid:  IsNotNil(args.CreateTime),
		},
		Name: name,
	})

	if err != nil {
		return fmt.Errorf("err updating env: %w", err)
	}

	return nil
}

func (e *EnvService) CreateEnvVar(ctx context.Context, args domain.CreateEnvVarArgs) (*domain.EnvVar, error) {
	return nil, errors.New("TODO")
}
