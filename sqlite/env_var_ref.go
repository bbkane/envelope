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
