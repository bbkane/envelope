package sqlite

import (
	"context"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

func (e *EnvService) EnvLocalVarCreate(ctx context.Context, args domain.EnvLocalVarCreateArgs) (*domain.EnvLocalVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.FindEnvID(ctx, args.EnvName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", args.Name, err)
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

func (e *EnvService) EnvLocalVarList(ctx context.Context, envName string) ([]domain.EnvLocalVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.FindEnvID(ctx, envName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", envName, err)
	}

	envs, err := queries.EnvLocalVarList(ctx, envID)
	if err != nil {
		return nil, fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}
	var ret []domain.EnvLocalVar
	for _, sqlcEnv := range envs {

		createTime, err := domain.StringToTime(sqlcEnv.CreateTime)
		if err != nil {
			return nil, fmt.Errorf("invalid create time for env_var %s: %w", sqlcEnv.Name, err)
		}

		updateTime, err := domain.StringToTime(sqlcEnv.UpdateTime)
		if err != nil {
			return nil, fmt.Errorf("invalid update time for env_var %s: %w", sqlcEnv.Name, err)
		}

		ret = append(ret, domain.EnvLocalVar{
			Name:       sqlcEnv.Name,
			Comment:    sqlcEnv.Comment,
			CreateTime: createTime,
			EnvName:    envName,
			UpdateTime: updateTime,
			Value:      sqlcEnv.Value,
		})
	}

	return ret, nil
}

func (e *EnvService) EnvLocalVarShow(ctx context.Context, envName string, name string) (*domain.EnvLocalVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.FindEnvID(ctx, envName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", envName, err)
	}

	sqlEnvLocalVar, err := queries.FindEnvLocalVar(ctx, sqlcgen.FindEnvLocalVarParams{
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
