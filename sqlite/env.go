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

	comment := sql.NullString{
		String: DerefOrEmpty(args.Comment),
		Valid:  IsNotNil(args.Comment),
	}

	createTime, err := domain.TimeToString(args.CreateTime)
	if err != nil {
		return 0, fmt.Errorf("could not translate createTime into string: %w", err)
	}

	updateTime, err := domain.TimeToString(args.UpdateTime)
	if err != nil {
		return 0, fmt.Errorf("could not translate UpdateTime into string: %w", err)
	}

	createdEnvID, err := queries.CreateEnv(ctx, sqlcgen.CreateEnvParams{
		Name:       args.Name,
		Comment:    comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
	})
	if err != nil {
		return 0, fmt.Errorf("could not create env in db: %w", err)
	}
	return domain.EnvID(createdEnvID), nil
}

func (e *EnvService) UpdateEnv(ctx context.Context, args domain.UpdateEnvArgs) error {
	// queries := sqlcgen.New(e.db)

	return errors.New("TODO")
}
