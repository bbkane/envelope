package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlite"
)

type EnvService struct {
	db *sql.DB
}

func NewEnvService(ctx context.Context, dsn string) (domain.EnvService, error) {
	db, err := sqlite.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}
	return &EnvService{
		db: db,
	}, nil
}
