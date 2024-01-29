package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

// keyringEntryCheckData checks whether the keyring exists in both the db, and the keychain
func (e *EnvService) keyringEntryCheckData(ctx context.Context, name string) (bool, bool, error) {

	// check if it's in the db
	queries := sqlcgen.New(e.db)
	existsInDB := false
	_, err := queries.KeyringEntryFindID(ctx, name)
	if err == nil {
		existsInDB = true
	} else if err != sql.ErrNoRows {
		return false, false, err
	}

	// check if exists in keyring
	existsInKeyring := false
	_, err = e.keyring.Get(name)

	if err == nil {
		existsInKeyring = true
	} else if err != domain.ErrNotFound {
		return false, false, err
	}

	return existsInDB, existsInKeyring, nil
}

func (e *EnvService) KeyringEntryCreate(ctx context.Context, args domain.KeyringEntryCreateArgs) (*domain.KeyringEntry, error) {

	existsInDB, existsInKeyring, err := e.keyringEntryCheckData(ctx, args.Name)

	if err != nil {
		return nil, err
	}

	if existsInDB && existsInKeyring {
		return nil, errors.New("data exists")

	} else if existsInDB && !existsInKeyring {
		return nil, errors.New("data exists in db, not in keychain")

	} else if !existsInDB && existsInKeyring {
		return nil, errors.New("data exist in keychain, not in db")

	} else {
		queries := sqlcgen.New(e.db)

		err = e.keyring.Set(args.Name, args.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set value in keyring: %w", err)
		}

		err = queries.KeyringEntryCreate(ctx, sqlcgen.KeyringEntryCreateParams{
			Name:       args.Name,
			Comment:    args.Comment,
			CreateTime: domain.TimeToString(args.CreateTime),
			UpdateTime: domain.TimeToString(args.UpdateTime),
		})
		if err != nil {
			return nil, fmt.Errorf("val in keyring, but not in db: (%s, %s) %w", e.keyring.Service(), args.Name, err)
		}
		return &domain.KeyringEntry{
			Name:       args.Name,
			Comment:    args.Comment,
			CreateTime: args.CreateTime,
			UpdateTime: args.UpdateTime,
			Value:      args.Value,
		}, nil
	}

}
