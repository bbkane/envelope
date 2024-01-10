package sqlite

import (
	"context"
	"database/sql"
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

func (e *EnvService) CreateEnv(ctx context.Context, args domain.CreateEnvArgs) (domain.EnvID, error) {
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
		return 0, fmt.Errorf("could not create env in db: %w", err)
	}
	return domain.EnvID(createdEnvID), nil
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
		Name: name, // TODO: fix
	})

	if err != nil {
		return fmt.Errorf("err updating env: %w", err)
	}

	return nil
}
