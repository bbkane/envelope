package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlite"
)

type EnvService struct {
	db      *sql.DB
	keyring domain.Keyring
}

func NewEnvService(ctx context.Context, dsn string, keyring domain.Keyring) (domain.EnvService, error) {
	// TODO use context!!
	db, err := sqlite.Connect(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}
	return &EnvService{
		db:      db,
		keyring: keyring,
	}, nil
}
