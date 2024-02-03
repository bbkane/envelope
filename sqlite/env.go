package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

func mapErrEnvNotFound(e error) error {
	if errors.Is(e, sql.ErrNoRows) {
		return domain.ErrEnvNotFound
	} else {
		return e
	}
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
		return nil, mapErrEnvNotFound(err)
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

	var createTimeStr *string
	if args.CreateTime != nil {
		tmp := domain.TimeToString(*args.CreateTime)
		createTimeStr = &tmp
	}

	var updateTimeStr *string
	if args.CreateTime != nil {
		tmp := domain.TimeToString(*args.UpdateTime)
		updateTimeStr = &tmp
	}

	err := queries.EnvUpdate(ctx, sqlcgen.EnvUpdateParams{
		NewName:    args.NewName,
		Comment:    args.Comment,
		CreateTime: createTimeStr,
		UpdateTime: updateTimeStr,
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
