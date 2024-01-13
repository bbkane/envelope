-- name: CreateEnv :one
INSERT INTO env (
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
)
RETURNING name, comment, create_time, update_time;

-- See https://docs.sqlc.dev/en/latest/howto/named_parameters.html#nullable-parameters
-- name: UpdateEnv :exec
UPDATE env SET
    name = COALESCE(sqlc.narg('new_name'), name),
    comment = COALESCE(sqlc.narg('comment'), comment),
    create_time = COALESCE(sqlc.narg('create_time'), create_time),
    update_time = COALESCE(sqlc.narg('update_time'), update_time)
WHERE name = sqlc.arg('name');

-- name: FindEnvID :one
SELECT id FROM env WHERE name = ?;

-- name: CreateEnvVar :exec
INSERT INTO env_var (
    env_id, name, comment, create_time, update_time, type, local_value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?   , ?
);

-- name: ListEnvVars :many
SELECT * FROM env_var
WHERE env_id = ?
ORDER BY name ASC