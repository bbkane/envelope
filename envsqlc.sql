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

-- name: CreateLocalEnvVar :exec
INSERT INTO env_var_local(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: ListLocalEnvVars :many
SELECT * FROM env_var_local
WHERE env_id = ?
ORDER BY name ASC;

-- name: CreateKeyringEntry :exec
INSERT INTO keyring_entry(
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
);