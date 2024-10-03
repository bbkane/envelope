-- name: VarRefCreate :exec
INSERT INTO env_ref(
    env_id, name, comment, create_time, update_time, env_var_id
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: VarRefDelete :exec
DELETE FROM env_ref WHERE env_id = ? AND name = ?;

-- name: VarRefList :many
SELECT * FROM env_ref
WHERE env_id = ?
ORDER BY name ASC;

-- name: VarRefShow :one
SELECT *
FROM env_ref
WHERE env_id = ? AND name = ?;

-- name: VarRefListByVarID :many
SELECT env.name AS env_name, env_ref.* FROM env_ref
JOIN env ON env_ref.env_id = env.env_id
WHERE env_var_id = ?
ORDER BY env_ref.name ASC;