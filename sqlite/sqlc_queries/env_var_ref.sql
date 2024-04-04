-- name: EnvRefCreate :exec
INSERT INTO env_ref(
    env_id, name, comment, create_time, update_time, env_var_id
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: EnvRefDelete :exec
DELETE FROM env_ref WHERE env_id = ? AND name = ?;

-- name: EnvRefList :many
SELECT * FROM env_ref
WHERE env_id = ?
ORDER BY name ASC;

-- name: EnvRefShow :one
SELECT *
FROM env_ref
WHERE env_id = ? AND name = ?;

-- name: EnvRefListByEnvVarID :many
SELECT env.name AS env_name, env_ref.* FROM env_ref
JOIN env ON env_ref.env_id = env.env_id
WHERE env_var_id = ?
ORDER BY env_ref.name ASC;