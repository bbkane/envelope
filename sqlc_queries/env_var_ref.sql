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