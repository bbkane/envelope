package sqlite

import (
	"context"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

func (e *EnvService) EnvLocalVarCreate(ctx context.Context, args domain.EnvLocalVarCreateArgs) (*domain.EnvLocalVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, args.EnvName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", args.Name, mapErrEnvNotFound(err))
	}

	err = queries.EnvLocalVarCreate(ctx, sqlcgen.EnvLocalVarCreateParams{
		EnvID:      envID,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: domain.TimeToString(args.CreateTime),
		UpdateTime: domain.TimeToString(args.UpdateTime),
		Value:      args.Value,
	})

	if err != nil {
		return nil, fmt.Errorf("could not create env var: %w", err)
	}
	return &domain.EnvLocalVar{
		EnvName:    args.EnvName,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: args.CreateTime,
		UpdateTime: args.UpdateTime,
		Value:      args.Value,
	}, nil
}

func (e *EnvService) EnvLocalVarDelete(ctx context.Context, envName string, name string) error {
	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, envName)
	if err != nil {
		return fmt.Errorf("could not find env with name: %s: %w", envName, mapErrEnvNotFound(err))
	}

	err = queries.EnvLocalVarDelete(ctx, sqlcgen.EnvLocalVarDeleteParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("could not delete env var: %s: %s: %w", envName, name, err)
	}
	return nil
}

func (e *EnvService) EnvLocalVarList(ctx context.Context, envName string) ([]domain.EnvLocalVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, envName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", envName, mapErrEnvNotFound(err))
	}

	envs, err := queries.EnvLocalVarList(ctx, envID)
	if err != nil {
		return nil, fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}
	var ret []domain.EnvLocalVar
	for _, sqlcEnv := range envs {
		ret = append(ret, domain.EnvLocalVar{
			Name:       sqlcEnv.Name,
			Comment:    sqlcEnv.Comment,
			CreateTime: domain.StringToTimeMust(sqlcEnv.CreateTime),
			EnvName:    envName,
			UpdateTime: domain.StringToTimeMust(sqlcEnv.UpdateTime),
			Value:      sqlcEnv.Value,
		})
	}

	return ret, nil
}

func (e *EnvService) EnvLocalVarShow(ctx context.Context, envName string, name string) (*domain.EnvLocalVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, envName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", envName, mapErrEnvNotFound(err))
	}

	sqlEnvLocalVar, err := queries.EnvLocalVarShow(ctx, sqlcgen.EnvLocalVarShowParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not find env var: %s: %s: %w", envName, name, err)
	}
	return &domain.EnvLocalVar{
		EnvName:    envName,
		Name:       name,
		Comment:    sqlEnvLocalVar.Comment,
		CreateTime: domain.StringToTimeMust(sqlEnvLocalVar.CreateTime),
		UpdateTime: domain.StringToTimeMust(sqlEnvLocalVar.UpdateTime),
		Value:      sqlEnvLocalVar.Value,
	}, nil
}
