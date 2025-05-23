package app

import (
	"context"
	"database/sql"
	"fmt"

	"go.bbkane.com/enventory/app/sqliteconnect"
	"go.bbkane.com/enventory/models"
)

type EnvService struct {
	db *sql.DB
}

func NewEnvService(ctx context.Context, dsn string) (models.EnvService, error) {
	db, err := sqliteconnect.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}
	return &EnvService{
		db: db,
	}, nil
}
