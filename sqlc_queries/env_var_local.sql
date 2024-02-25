-- name: EnvVarCreate :exec
INSERT INTO env_var_local(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: EnvVarFindID :one
SELECT id FROM env_var_local WHERE env_id = ? AND name = ?;

-- name: EnvVarFindByID :one
SELECT * FROM env_var_local WHERE id = ?;

-- name: EnvVarDelete :exec
DELETE FROM env_var_local WHERE env_id = ? AND name = ?;

-- name: EnvVarList :many
SELECT * FROM env_var_local
WHERE env_id = ?
ORDER BY name ASC;

-- name: EnvVarShow :one
SELECT *
FROM env_var_local
WHERE env_id = ? AND name = ?;