// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: keyring_entry.sql

package sqlcgen

import (
	"context"
)

const keyringEntryCreate = `-- name: KeyringEntryCreate :exec
INSERT INTO keyring_entry(
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
)
`

type KeyringEntryCreateParams struct {
	Name       string
	Comment    string
	CreateTime string
	UpdateTime string
}

func (q *Queries) KeyringEntryCreate(ctx context.Context, arg KeyringEntryCreateParams) error {
	_, err := q.db.ExecContext(ctx, keyringEntryCreate,
		arg.Name,
		arg.Comment,
		arg.CreateTime,
		arg.UpdateTime,
	)
	return err
}

const keyringEntryDelete = `-- name: KeyringEntryDelete :exec
DELETE FROM keyring_entry WHERE name = ?
`

func (q *Queries) KeyringEntryDelete(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, keyringEntryDelete, name)
	return err
}

const keyringEntryFindID = `-- name: KeyringEntryFindID :one
SELECT keyring_entry_id from keyring_entry WHERE name = ?
`

func (q *Queries) KeyringEntryFindID(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, keyringEntryFindID, name)
	var keyring_entry_id int64
	err := row.Scan(&keyring_entry_id)
	return keyring_entry_id, err
}

const keyringEntryList = `-- name: KeyringEntryList :many
SELECT keyring_entry_id, name, comment, create_time, update_time FROM keyring_entry
ORDER BY name ASC
`

func (q *Queries) KeyringEntryList(ctx context.Context) ([]KeyringEntry, error) {
	rows, err := q.db.QueryContext(ctx, keyringEntryList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []KeyringEntry
	for rows.Next() {
		var i KeyringEntry
		if err := rows.Scan(
			&i.KeyringEntryID,
			&i.Name,
			&i.Comment,
			&i.CreateTime,
			&i.UpdateTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const keyringEntryShow = `-- name: KeyringEntryShow :one
SELECT keyring_entry_id, name, comment, create_time, update_time FROM keyring_entry WHERE name = ?
`

func (q *Queries) KeyringEntryShow(ctx context.Context, name string) (KeyringEntry, error) {
	row := q.db.QueryRowContext(ctx, keyringEntryShow, name)
	var i KeyringEntry
	err := row.Scan(
		&i.KeyringEntryID,
		&i.Name,
		&i.Comment,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}
