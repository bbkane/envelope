package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlcgen"
)

// mapErrEnvNotFound replaces sql.ErrNoRows with domain.ErrEnvNotFound but otherwise
// passes it through.
//
// Deprecated: I want to replace this with envFindID, but that'll require rewriting some sql
func mapErrEnvNotFound(e error) error {
	if errors.Is(e, sql.ErrNoRows) {
		return domain.ErrEnvNotFound
	} else {
		return e
	}
}

// envFindID looks for en env's SQLite ID and returns a wrapped ErrEnvNotFound or a sql error
func (e *EnvService) envFindID(ctx context.Context, envName string) (int64, error) {
	queries := sqlcgen.New(e.db)
	envID, err := queries.EnvFindID(ctx, envName)
	if errors.Is(err, sql.ErrNoRows) {
		err = domain.ErrEnvNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("could not find env with name: %s: %w", envName, err)
	}
	return envID, nil
}

func (e *EnvService) EnvCreate(ctx context.Context, args domain.EnvCreateArgs) (*domain.Env, error) {
	queries := sqlcgen.New(e.db)

	createdEnvID, err := queries.EnvCreate(ctx, sqlcgen.EnvCreateParams{
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: domain.TimeToString(args.CreateTime),
		UpdateTime: domain.TimeToString(args.UpdateTime),
	})

	if err != nil {
		return nil, fmt.Errorf("could not create env in db: %w", err)
	}

	return &domain.Env{
		Name:       createdEnvID.Name,
		Comment:    createdEnvID.Comment,
		CreateTime: domain.StringToTimeMust(createdEnvID.CreateTime),
		UpdateTime: domain.StringToTimeMust(createdEnvID.UpdateTime),
	}, nil
}

func (e *EnvService) EnvDelete(ctx context.Context, name string) error {
	queries := sqlcgen.New(e.db)

	err := queries.EnvDelete(ctx, name)
	if err != nil {
		return mapErrEnvNotFound(err)
	}
	return nil
}

func (e *EnvService) EnvList(ctx context.Context) ([]domain.Env, error) {
	queries := sqlcgen.New(e.db)

	sqlcEnvs, err := queries.EnvList(ctx)
	if err != nil {
		return nil, err
	}

	ret := []domain.Env{}
	for _, e := range sqlcEnvs {
		ret = append(ret, domain.Env{
			Name:       e.Name,
			Comment:    e.Comment,
			CreateTime: domain.StringToTimeMust(e.CreateTime),
			UpdateTime: domain.StringToTimeMust(e.UpdateTime),
		})
	}

	return ret, nil
}

func (e *EnvService) EnvUpdate(ctx context.Context, name string, args domain.EnvUpdateArgs) error {

	queries := sqlcgen.New(e.db)

	err := queries.EnvUpdate(ctx, sqlcgen.EnvUpdateParams{
		NewName:    args.Name,
		Comment:    args.Comment,
		CreateTime: domain.TimePtrToStringPtr(args.CreateTime),
		UpdateTime: domain.TimePtrToStringPtr(args.UpdateTime),
		Name:       name,
	})

	if err != nil {
		return fmt.Errorf("err updating env: %w", mapErrEnvNotFound(err))
	}

	return nil
}

func (e *EnvService) EnvShow(ctx context.Context, name string) (*domain.Env, error) {
	queries := sqlcgen.New(e.db)

	sqlcEnv, err := queries.EnvShow(ctx, name)

	if err != nil {
		return nil, mapErrEnvNotFound(err)
	}

	return &domain.Env{
		Name:       name,
		Comment:    sqlcEnv.Comment,
		CreateTime: domain.StringToTimeMust(sqlcEnv.CreateTime),
		UpdateTime: domain.StringToTimeMust(sqlcEnv.UpdateTime),
	}, nil
}
