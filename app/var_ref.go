package app

import (
	"context"
	"fmt"

	"go.bbkane.com/envelope/app/sqliteconnect/sqlcgen"
	"go.bbkane.com/envelope/models"
)

func (e *EnvService) VarRefCreate(ctx context.Context, args models.VarRefCreateArgs) (*models.VarRef, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, args.EnvName)
	if err != nil {
		return nil, err
	}

	varID, err := e.varFindID(ctx, args.RefEnvName, args.RefVarName)
	if err != nil {
		return nil, err
	}

	err = queries.VarRefCreate(ctx, sqlcgen.VarRefCreateParams{
		EnvID:      envID,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: models.TimeToString(args.CreateTime),
		UpdateTime: models.TimeToString(args.UpdateTime),
		VarID:      varID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create env var ref: %w", err)
	}
	return &models.VarRef{
		EnvName:    args.EnvName,
		Name:       args.Name,
		Comment:    args.Comment,
		CreateTime: args.CreateTime,
		UpdateTime: args.UpdateTime,
		RefEnvName: args.RefEnvName,
		RevVarName: args.RefVarName,
	}, nil
}

func (e *EnvService) VarRefDelete(ctx context.Context, envName string, name string) error {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return err
	}

	rowsAffected, err := queries.VarRefDelete(ctx, sqlcgen.VarRefDeleteParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return fmt.Errorf("could not delete ref: %s: %s: %w", envName, name, err)
	}
	if rowsAffected == 0 {
		return models.ErrVarRefNotFound
	}
	return nil
}

func (e *EnvService) VarRefList(ctx context.Context, envName string) ([]models.VarRef, []models.Var, error) {
	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, nil, err
	}

	sqlcRefs, err := queries.VarRefList(ctx, envID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}
	var refs []models.VarRef
	var vars []models.Var
	for _, sqlcRef := range sqlcRefs {

		// classic N+1 query pattern, but luckily SQLite is not really affected by this :)
		// https://www.sqlite.org/np1queryprob.html
		// easy to add a join later if I need perf, as this is localized to this package

		localVar, err := e.varFindByID(ctx, sqlcRef.VarID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not find var from id: %d: %w", sqlcRef.VarID, err)
		}
		vars = append(vars, *localVar)
		refs = append(refs, models.VarRef{
			EnvName:    envName,
			Name:       sqlcRef.Name,
			Comment:    sqlcRef.Comment,
			CreateTime: models.StringToTimeMust(sqlcRef.CreateTime),
			UpdateTime: models.StringToTimeMust(sqlcRef.UpdateTime),
			RefEnvName: localVar.EnvName,
			RevVarName: localVar.Name,
		})
	}

	return refs, vars, nil
}

func (e *EnvService) VarRefShow(ctx context.Context, envName string, name string) (*models.VarRef, *models.Var, error) {

	queries := sqlcgen.New(e.db)

	envID, err := e.envFindID(ctx, envName)
	if err != nil {
		return nil, nil, err
	}

	sqlcRef, err := queries.VarRefShow(ctx, sqlcgen.VarRefShowParams{
		EnvID: envID,
		Name:  name,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not find ref: %s: %s: %w", envName, name, err)
	}
	sqlcVar, err := e.varFindByID(ctx, sqlcRef.VarID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find var from id: %d: %w", sqlcRef.VarID, err)
	}

	return &models.VarRef{
			EnvName:    envName,
			Name:       sqlcRef.Name,
			Comment:    sqlcRef.Comment,
			CreateTime: models.StringToTimeMust(sqlcRef.CreateTime),
			UpdateTime: models.StringToTimeMust(sqlcRef.UpdateTime),
			RefEnvName: sqlcVar.EnvName,
			RevVarName: sqlcVar.Name,
		}, &models.Var{
			EnvName:    sqlcVar.EnvName,
			Name:       sqlcVar.Name,
			Comment:    sqlcVar.Comment,
			CreateTime: sqlcVar.CreateTime,
			UpdateTime: sqlcVar.UpdateTime,
			Value:      sqlcVar.Value,
		}, nil
}
