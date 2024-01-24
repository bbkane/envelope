// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: env_var_local.sql

package sqlcgen

import (
	"context"
)

const envLocalVarCreate = `-- name: EnvLocalVarCreate :exec
INSERT INTO env_var_local(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
)
`

type EnvLocalVarCreateParams struct {
	EnvID      int64
	Name       string
	Comment    string
	CreateTime string
	UpdateTime string
	Value      string
}

func (q *Queries) EnvLocalVarCreate(ctx context.Context, arg EnvLocalVarCreateParams) error {
	_, err := q.db.ExecContext(ctx, envLocalVarCreate,
		arg.EnvID,
		arg.Name,
		arg.Comment,
		arg.CreateTime,
		arg.UpdateTime,
		arg.Value,
	)
	return err
}

const envLocalVarDelete = `-- name: EnvLocalVarDelete :exec
DELETE FROM env_var_local WHERE env_id = ? AND  name = ?
`

type EnvLocalVarDeleteParams struct {
	EnvID int64
	Name  string
}

func (q *Queries) EnvLocalVarDelete(ctx context.Context, arg EnvLocalVarDeleteParams) error {
	_, err := q.db.ExecContext(ctx, envLocalVarDelete, arg.EnvID, arg.Name)
	return err
}

const envLocalVarList = `-- name: EnvLocalVarList :many
SELECT id, env_id, name, comment, create_time, update_time, value FROM env_var_local
WHERE env_id = ?
ORDER BY name ASC
`

func (q *Queries) EnvLocalVarList(ctx context.Context, envID int64) ([]EnvVarLocal, error) {
	rows, err := q.db.QueryContext(ctx, envLocalVarList, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EnvVarLocal
	for rows.Next() {
		var i EnvVarLocal
		if err := rows.Scan(
			&i.ID,
			&i.EnvID,
			&i.Name,
			&i.Comment,
			&i.CreateTime,
			&i.UpdateTime,
			&i.Value,
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

const findEnvLocalVar = `-- name: FindEnvLocalVar :one
SELECT id, env_id, name, comment, create_time, update_time, value
FROM env_var_local
WHERE env_id = ? AND name = ?
`

type FindEnvLocalVarParams struct {
	EnvID int64
	Name  string
}

func (q *Queries) FindEnvLocalVar(ctx context.Context, arg FindEnvLocalVarParams) (EnvVarLocal, error) {
	row := q.db.QueryRowContext(ctx, findEnvLocalVar, arg.EnvID, arg.Name)
	var i EnvVarLocal
	err := row.Scan(
		&i.ID,
		&i.EnvID,
		&i.Name,
		&i.Comment,
		&i.CreateTime,
		&i.UpdateTime,
		&i.Value,
	)
	return i, err
}
