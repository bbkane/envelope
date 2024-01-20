package sqlite

import (
	"context"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

func (e *EnvService) EnvCreate(ctx context.Context, args domain.EnvCreateArgs) (*domain.Env, error) {
	queries := sqlcgen.New(e.db)

	createdEnvID, err := queries.CreateEnv(ctx, sqlcgen.CreateEnvParams{
		Name:       args.Name,
		Comment:    args.Comment,
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
		Comment:    createdEnvID.Comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}, nil
}

func (e *EnvService) EnvDelete(ctx context.Context, name string) error {
	queries := sqlcgen.New(e.db)

	err := queries.DeleteEnv(ctx, name)
	if err != nil {
		return err
	}
	return nil
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

	err := queries.UpdateEnv(ctx, sqlcgen.UpdateEnvParams{
		NewName:    args.NewName,
		Comment:    args.Comment,
		CreateTime: createTimeStr,
		UpdateTime: updateTimeStr,
		Name:       name,
	})

	if err != nil {
		return fmt.Errorf("err updating env: %w", err)
	}

	return nil
}

func (e *EnvService) EnvShow(ctx context.Context, name string) (*domain.Env, error) {
	queries := sqlcgen.New(e.db)

	sqlcEnv, err := queries.FindEnv(ctx, name)

	if err != nil {
		return nil, fmt.Errorf("could not find env: %s: %w", name, err)
	}

	createTime, err := domain.StringToTime(sqlcEnv.CreateTime)
	if err != nil {
		return nil, fmt.Errorf("bad create_time: %s: %w", name, err)
	}

	updateTime, err := domain.StringToTime(sqlcEnv.UpdateTime)
	if err != nil {
		return nil, fmt.Errorf("bad update_time: %s: %w", name, err)
	}

	return &domain.Env{
		Name:       name,
		Comment:    sqlcEnv.Comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}, nil
}
