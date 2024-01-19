package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"
)

func (e *EnvService) KeyringEntryCreate(ctx context.Context, args domain.KeyringEntryCreateArgs) (*domain.KeyringEntry, error) {

	queries := sqlcgen.New(e.db)

	// check if the value exists before trying to create it
	// We assume if the keyring entry is in the db, it's also in the os keyring
	_, err := queries.FindKeyringID(ctx, args.Name)

	// oops, one exists
	if err == nil {
		return nil, errors.New("expecting no entries, but we found one")
	}

	// we expect this error, but we want to alert on any others
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("keyring query error: %w", err)
	}

	err = e.keyring.Set(args.Name, args.Value)
	if err != nil {
		return nil, fmt.Errorf("could not set value in keyring: %w", err)
	}

	err = queries.CreateKeyringEntry(ctx, sqlcgen.CreateKeyringEntryParams{
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
