-- name: EnvLocalVarCreate :exec
INSERT INTO env_var_local(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: EnvLocalVarFindID :one
SELECT id FROM env_var_local WHERE env_id = ? AND name = ?;

-- name: EnvLocalVarFindByID :one
SELECT * FROM env_var_local WHERE id = ?;

-- name: EnvLocalVarDelete :exec
DELETE FROM env_var_local WHERE env_id = ? AND name = ?;

-- name: EnvLocalVarList :many
SELECT * FROM env_var_local
WHERE env_id = ?
ORDER BY name ASC;

-- name: EnvLocalVarShow :one
SELECT *
FROM env_var_local
WHERE env_id = ? AND name = ?;