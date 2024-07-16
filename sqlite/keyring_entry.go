package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/sqlite/sqlite/sqlcgen"
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
	} else if err != domain.ErrKeyringSecretNotFound {
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

// KeyringEntryLists returns a list of KeyringEntrys, a list of errors when retrieving individual values from the keyring, and an error if the whole operation fails
func (e *EnvService) KeyringEntryList(ctx context.Context) ([]domain.KeyringEntry, []error, error) {

	queries := sqlcgen.New(e.db)

	sqlcKEs, err := queries.KeyringEntryList(ctx)
	if err != nil {
		return nil, nil, err
	}
	errs := []error{}
	ret := []domain.KeyringEntry{}

	for _, sqlcKE := range sqlcKEs {
		// try to find it in the keyring
		kVal, err := e.keyring.Get(sqlcKE.Name)
		if err != nil {
			errs = append(errs, fmt.Errorf("error retriving from keyring: %s: %w", sqlcKE.Name, err))
		}
		ret = append(ret, domain.KeyringEntry{
			Name:       sqlcKE.Name,
			Comment:    sqlcKE.Comment,
			CreateTime: domain.StringToTimeMust(sqlcKE.CreateTime),
			UpdateTime: domain.StringToTimeMust(sqlcKE.UpdateTime),
			Value:      kVal,
		})
	}
	return ret, errs, err
}
