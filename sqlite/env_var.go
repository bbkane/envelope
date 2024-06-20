package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlcgen"
)

func (e *EnvService) envLocalVarFindByID(ctx context.Context, id int64) (*domain.EnvVar, error) {
	queries := sqlcgen.New(e.db)

	sqlcVar, err := queries.EnvVarFindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrEnvVarNotFound
	}

	return &domain.EnvVar{
		EnvName:    sqlcVar.EnvName,
		Name:       sqlcVar.Name,
		Comment:    sqlcVar.Comment,
		CreateTime: domain.StringToTimeMust(sqlcVar.CreateTime),
		UpdateTime: domain.StringToTimeMust(sqlcVar.UpdateTime),
		Value:      sqlcVar.Value,
	}, nil
}

func (e *EnvService) envLocalVarFindID(ctx context.Context, envName string, name string) (int64, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return 0, err
	}

	id, err := queries.EnvVarFindID(ctx, sqlcgen.EnvVarFindIDParams{
		EnvID: envID,
		Name:  name,
	})

	if err != nil {
		return 0, domain.ErrEnvVarNotFound
	}
	return id, nil

}

func (e *EnvService) EnvVarCreate(ctx context.Context, args domain.EnvVarCreateArgs) (*domain.EnvVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, args.EnvName)
	if err != nil {
		return nil, err
	}

	err = queries.EnvVarCreate(ctx, sqlcgen.EnvVarCreateParams{
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
	return &domain.EnvVar{
		EnvName:    args.EnvName,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: args.CreateTime,
		UpdateTime: args.UpdateTime,
		Value:      args.Value,
	}, nil
}

func (e *EnvService) EnvVarDelete(ctx context.Context, envName string, name string) error {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return err
	}

	err = queries.EnvVarDelete(ctx, sqlcgen.EnvVarDeleteParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("could not delete env var: %s: %s: %w", envName, name, err)
	}
	return nil
}

func (e *EnvService) EnvVarList(ctx context.Context, envName string) ([]domain.EnvVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, err
	}

	envs, err := queries.EnvVarList(ctx, envID)
	if err != nil {
		return nil, fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}
	var ret []domain.EnvVar
	for _, sqlcEnv := range envs {
		ret = append(ret, domain.EnvVar{
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

func (e *EnvService) EnvVarShow(ctx context.Context, envName string, name string) (*domain.EnvVar, []domain.EnvRef, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, nil, err
	}

	sqlEnvLocalVar, err := queries.EnvVarShow(ctx, sqlcgen.EnvVarShowParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not find env var: %s: %s: %w", envName, name, err)
	}

	envRefs := []domain.EnvRef{}
	sqlcEnvRefs, err := queries.EnvRefListByEnvVarID(ctx, sqlEnvLocalVar.EnvVarID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}

	for _, e := range sqlcEnvRefs {
		envRefs = append(envRefs, domain.EnvRef{
			EnvName:    e.EnvName,
			Name:       e.Name,
			Comment:    e.Comment,
			CreateTime: domain.StringToTimeMust(e.CreateTime),
			UpdateTime: domain.StringToTimeMust(e.UpdateTime),
			RefEnvName: envName,
			RevVarName: name,
		})
	}

	return &domain.EnvVar{
		EnvName:    envName,
		Name:       name,
		Comment:    sqlEnvLocalVar.Comment,
		CreateTime: domain.StringToTimeMust(sqlEnvLocalVar.CreateTime),
		UpdateTime: domain.StringToTimeMust(sqlEnvLocalVar.UpdateTime),
		Value:      sqlEnvLocalVar.Value,
	}, envRefs, nil
}

func (e *EnvService) EnvVarUpdate(ctx context.Context, envName string, name string, args domain.EnvVarUpdateArgs) error {
	return errors.New("TODO")
}
