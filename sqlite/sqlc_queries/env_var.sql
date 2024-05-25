-- name: EnvVarCreate :exec
INSERT INTO env_var(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: EnvVarFindID :one
SELECT env_var_id FROM env_var WHERE env_id = ? AND name = ?;

-- name: EnvVarFindByID :one
SELECT env.name AS env_name, env_var.*
FROM env_var
JOIN env ON env_var.env_id = env.env_id
WHERE env_var.env_var_id = ?;

-- name: EnvVarDelete :exec
DELETE FROM env_var WHERE env_id = ? AND name = ?;

-- name: EnvVarList :many
SELECT * FROM env_var
WHERE env_id = ?
ORDER BY name ASC;

-- name: EnvVarShow :one
SELECT *
FROM env_var
WHERE env_id = ? AND name = ?;