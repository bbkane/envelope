-- name: EnvVarCreate :exec
INSERT INTO env_var(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: EnvVarDelete :exec
DELETE FROM env_var WHERE env_id = ? AND name = ?;

-- name: EnvVarFindID :one
SELECT env_var_id FROM env_var WHERE env_id = ? AND name = ?;

-- name: EnvVarFindByID :one
SELECT env.name AS env_name, env_var.*
FROM env_var
JOIN env ON env_var.env_id = env.env_id
WHERE env_var.env_var_id = ?;

-- name: EnvVarList :many
SELECT * FROM env_var
WHERE env_id = ?
ORDER BY name ASC;

-- name: EnvVarShow :one
SELECT *
FROM env_var
WHERE env_id = ? AND name = ?;

-- name: EnvVarUpdate :exec
UPDATE env_var SET
    env_id = COALESCE(sqlc.narg('env_id'), env_id),
    name = COALESCE(sqlc.narg('name'), name),
    comment = COALESCE(sqlc.narg('comment'), comment),
    create_time = COALESCE(sqlc.narg('create_time'), create_time),
    update_time = COALESCE(sqlc.narg('update_time'), update_time),
    value = COALESCE(sqlc.narg('value'), value)
WHERE env_var_id = sqlc.arg('env_var_id');