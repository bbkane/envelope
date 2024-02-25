package sqlite

import (
	"context"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlcgen"
)

func (e *EnvService) EnvRefCreate(ctx context.Context, args domain.EnvLocalRefCreateArgs) (*domain.EnvLocalRef, error) {
	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, args.EnvName)
	if err != nil {
		return nil, fmt.Errorf("could not find env with name: %s: %w", args.Name, mapErrEnvNotFound(err))
	}

	varID, err := e.envLocalVarFindID(ctx, args.RefEnvName, args.RefVarName)
	if err != nil {
		return nil, err
	}

	err = queries.EnvRefCreate(ctx, sqlcgen.EnvRefCreateParams{
		EnvID:         envID,
		Name:          args.Name,
		Comment:       args.Comment,
		CreateTime:    domain.TimeToString(args.CreateTime),
		UpdateTime:    domain.TimeToString(args.UpdateTime),
		EnvVarLocalID: varID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create env var ref: %w", err)
	}
	return &domain.EnvLocalRef{
		EnvName:    args.EnvName,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: args.CreateTime,
		UpdateTime: args.UpdateTime,
		RefEnvName: args.RefEnvName,
		RevVarName: args.RefVarName,
	}, nil
}

func (e *EnvService) EnvRefDelete(ctx context.Context, envName string, name string) error {
	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, envName)
	if err != nil {
		return fmt.Errorf("could not find env with name: %s: %w", envName, mapErrEnvNotFound(err))
	}

	err = queries.EnvRefDelete(ctx, sqlcgen.EnvRefDeleteParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("could not delete ref: %s: %s: %w", envName, name, err)
	}
	return nil
}

func (e *EnvService) EnvRefList(ctx context.Context, envName string) ([]domain.EnvLocalRef, []domain.EnvLocalVar, error) {
	panic("TODO")
}

func (e *EnvService) EnvRefShow(ctx context.Context, envName string, name string) (*domain.EnvLocalRef, *domain.EnvLocalVar, error) {

	queries := sqlcgen.New(e.db)

	envID, err := queries.EnvFindID(ctx, envName)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find env with name: %s: %w", envName, mapErrEnvNotFound(err))
	}

	sqlEnvRef, err := queries.EnvRefShow(ctx, sqlcgen.EnvRefShowParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not find ref: %s: %s: %w", envName, name, err)
	}
	sqlEnvVar, err := e.envLocalVarFindByID(ctx, envName, sqlEnvRef.EnvVarLocalID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find var from id: %d: %w", sqlEnvRef.EnvVarLocalID, err)
	}

	return &domain.EnvLocalRef{
			EnvName:    envName,
			Name:       sqlEnvRef.Name,
			Comment:    sqlEnvRef.Comment,
			CreateTime: domain.StringToTimeMust(sqlEnvRef.CreateTime),
			UpdateTime: domain.StringToTimeMust(sqlEnvRef.UpdateTime),
			RefEnvName: sqlEnvVar.EnvName,
			RevVarName: sqlEnvVar.Name,
		}, &domain.EnvLocalVar{
			EnvName:    sqlEnvVar.EnvName,
			Name:       sqlEnvVar.Name,
			Comment:    sqlEnvVar.Comment,
			CreateTime: sqlEnvVar.CreateTime,
			UpdateTime: sqlEnvVar.UpdateTime,
			Value:      sqlEnvVar.Value,
		}, nil
}
