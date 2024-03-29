-- name: EnvRefCreate :exec
INSERT INTO env_var_ref(
    env_id, name, comment, create_time, update_time, env_var_local_id
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: EnvRefDelete :exec
DELETE FROM env_var_ref WHERE env_id = ? AND name = ?;

-- name: EnvRefList :many
SELECT * FROM env_var_ref
WHERE env_id = ?
ORDER BY name ASC;

-- name: EnvRefShow :one
SELECT *
FROM env_var_ref
WHERE env_id = ? AND name = ?;

-- name: EnvRefListByEnvVarID :many
SELECT env.name AS env_name, env_var_ref.* FROM env_var_ref
JOIN env ON env_var_ref.env_id = env.id
WHERE env_var_local_id = ?
ORDER BY env_var_ref.name ASC;