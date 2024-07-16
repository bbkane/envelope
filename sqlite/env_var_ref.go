package sqlite

import (
	"context"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlite/sqlcgen"
)

func (e *EnvService) EnvRefCreate(ctx context.Context, args domain.EnvRefCreateArgs) (*domain.EnvRef, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, args.EnvName)
	if err != nil {
		return nil, err
	}

	varID, err := e.envLocalVarFindID(ctx, args.RefEnvName, args.RefVarName)
	if err != nil {
		return nil, err
	}

	err = queries.EnvRefCreate(ctx, sqlcgen.EnvRefCreateParams{
		EnvID:      envID,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: domain.TimeToString(args.CreateTime),
		UpdateTime: domain.TimeToString(args.UpdateTime),
		EnvVarID:   varID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create env var ref: %w", err)
	}
	return &domain.EnvRef{
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

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return err
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

func (e *EnvService) EnvRefList(ctx context.Context, envName string) ([]domain.EnvRef, []domain.EnvVar, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, nil, err
	}

	sqlcRefs, err := queries.EnvRefList(ctx, envID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}
	var refs []domain.EnvRef
	var vars []domain.EnvVar
	for _, sqlcRef := range sqlcRefs {

		// classic N+1 query pattern, but luckily SQLite is not really affected by this :)
		// https://www.sqlite.org/np1queryprob.html
		// easy to add a join later if I need perf, as this is localized to this package

		localVar, err := e.envLocalVarFindByID(ctx, sqlcRef.EnvVarID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not find var from id: %d: %w", sqlcRef.EnvVarID, err)
		}
		vars = append(vars, *localVar)
		refs = append(refs, domain.EnvRef{
			EnvName:    envName,
			Name:       sqlcRef.Name,
			Comment:    sqlcRef.Comment,
			CreateTime: domain.StringToTimeMust(sqlcRef.CreateTime),
			UpdateTime: domain.StringToTimeMust(sqlcRef.UpdateTime),
			RefEnvName: localVar.EnvName,
			RevVarName: localVar.Name,
		})
	}

	return refs, vars, nil
}

func (e *EnvService) EnvRefShow(ctx context.Context, envName string, name string) (*domain.EnvRef, *domain.EnvVar, error) {

	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, nil, err
	}

	sqlcRef, err := queries.EnvRefShow(ctx, sqlcgen.EnvRefShowParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not find ref: %s: %s: %w", envName, name, err)
	}
	sqlcVar, err := e.envLocalVarFindByID(ctx, sqlcRef.EnvVarID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find var from id: %d: %w", sqlcRef.EnvVarID, err)
	}

	return &domain.EnvRef{
			EnvName:    envName,
			Name:       sqlcRef.Name,
			Comment:    sqlcRef.Comment,
			CreateTime: domain.StringToTimeMust(sqlcRef.CreateTime),
			UpdateTime: domain.StringToTimeMust(sqlcRef.UpdateTime),
			RefEnvName: sqlcVar.EnvName,
			RevVarName: sqlcVar.Name,
		}, &domain.EnvVar{
			EnvName:    sqlcVar.EnvName,
			Name:       sqlcVar.Name,
			Comment:    sqlcVar.Comment,
			CreateTime: sqlcVar.CreateTime,
			UpdateTime: sqlcVar.UpdateTime,
			Value:      sqlcVar.Value,
		}, nil
}
