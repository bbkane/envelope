package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/connect"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

type EnvService struct {
	db *sql.DB
}

func NewEnvService(ctx context.Context, dsn string) (domain.EnvService, error) {
	// TODO use context!!
	db, err := connect.Connect(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}
	return &EnvService{
		db: db,
	}, nil
}

func (e *EnvService) CreateEnv(ctx context.Context, args domain.CreateEnvArgs) (domain.EnvID, error) {
	queries := sqlcgen.New(e.db)

	var comment sql.NullString
	if args.Comment != nil {
		comment.String = *args.Comment
		comment.Valid = true
	}

	createdEnvID, err := queries.CreateEnv(ctx, sqlcgen.CreateEnvParams{
		Name:       args.Name,
		Comment:    comment,
		CreateTime: args.CreateTime.String(), // TODO: fix
		UpdateTime: args.UpdateTime.String(), // TODO: fix
	})
	if err != nil {
		return 0, fmt.Errorf("could not create env in db: %w", err)
	}
	return domain.EnvID(createdEnvID), nil
}
