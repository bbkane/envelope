package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlite/sqlcgen"
)

func (e *EnvService) varFindByID(ctx context.Context, id int64) (*domain.Var, error) {
	queries := sqlcgen.New(e.db)

	sqlcVar, err := queries.VarFindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrEnvVarNotFound
	}

	return &domain.Var{
		EnvName:    sqlcVar.EnvName,
		Name:       sqlcVar.Name,
		Comment:    sqlcVar.Comment,
		CreateTime: domain.StringToTimeMust(sqlcVar.CreateTime),
		UpdateTime: domain.StringToTimeMust(sqlcVar.UpdateTime),
		Value:      sqlcVar.Value,
	}, nil
}

func (e *EnvService) varFindID(ctx context.Context, envName string, name string) (int64, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return 0, err
	}

	id, err := queries.VarFindID(ctx, sqlcgen.VarFindIDParams{
		EnvID: envID,
		Name:  name,
	})

	if err != nil {
		return 0, domain.ErrEnvVarNotFound
	}
	return id, nil

}

func (e *EnvService) VarCreate(ctx context.Context, args domain.VarCreateArgs) (*domain.Var, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, args.EnvName)
	if err != nil {
		return nil, err
	}

	err = queries.VarCreate(ctx, sqlcgen.VarCreateParams{
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
	return &domain.Var{
		EnvName:    args.EnvName,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: args.CreateTime,
		UpdateTime: args.UpdateTime,
		Value:      args.Value,
	}, nil
}

func (e *EnvService) VarDelete(ctx context.Context, envName string, name string) error {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return err
	}

	err = queries.VarDelete(ctx, sqlcgen.VarDeleteParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("could not delete env var: %s: %s: %w", envName, name, err)
	}
	return nil
}

func (e *EnvService) VarList(ctx context.Context, envName string) ([]domain.Var, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, err
	}

	envs, err := queries.VarList(ctx, envID)
	if err != nil {
		return nil, fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}
	var ret []domain.Var
	for _, sqlcEnv := range envs {
		ret = append(ret, domain.Var{
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

func (e *EnvService) VarShow(ctx context.Context, envName string, name string) (*domain.Var, []domain.VarRef, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, nil, err
	}

	sqlEnvLocalVar, err := queries.VarShow(ctx, sqlcgen.VarShowParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not find env var: %s: %s: %w", envName, name, err)
	}

	envRefs := []domain.VarRef{}
	sqlcEnvRefs, err := queries.VarRefListByVarID(ctx, sqlEnvLocalVar.VarID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}

	for _, e := range sqlcEnvRefs {
		envRefs = append(envRefs, domain.VarRef{
			EnvName:    e.EnvName,
			Name:       e.Name,
			Comment:    e.Comment,
			CreateTime: domain.StringToTimeMust(e.CreateTime),
			UpdateTime: domain.StringToTimeMust(e.UpdateTime),
			RefEnvName: envName,
			RevVarName: name,
		})
	}

	return &domain.Var{
		EnvName:    envName,
		Name:       name,
		Comment:    sqlEnvLocalVar.Comment,
		CreateTime: domain.StringToTimeMust(sqlEnvLocalVar.CreateTime),
		UpdateTime: domain.StringToTimeMust(sqlEnvLocalVar.UpdateTime),
		Value:      sqlEnvLocalVar.Value,
	}, envRefs, nil
}

func (e *EnvService) VarUpdate(ctx context.Context, envName string, name string, args domain.VarUpdateArgs) error {
	envVarID, err := e.varFindID(ctx, envName, name)
	if err != nil {
		return err
	}

	var newEnvID *int64
	if args.EnvName != nil {
		tmp, err := e.envFindID(ctx, *args.EnvName)
		if err != nil {
			return err
		}
		newEnvID = &tmp
	}

	queries := sqlcgen.New(e.db)

	err = queries.VarUpdate(ctx, sqlcgen.VarUpdateParams{
		EnvID:      newEnvID,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: domain.TimePtrToStringPtr(args.CreateTime),
		UpdateTime: domain.TimePtrToStringPtr(args.UpdateTime),
		Value:      args.Value,
		VarID:      envVarID,
	})

	if err != nil {
		return fmt.Errorf("err updating env var: %w", err)
	}
	return nil
}
