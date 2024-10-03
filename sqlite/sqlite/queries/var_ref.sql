-- name: VarRefCreate :exec
INSERT INTO var_ref(
    env_id, name, comment, create_time, update_time, var_id
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: VarRefDelete :exec
DELETE FROM var_ref WHERE env_id = ? AND name = ?;

-- name: VarRefList :many
SELECT * FROM var_ref
WHERE env_id = ?
ORDER BY name ASC;

-- name: VarRefShow :one
SELECT *
FROM var_ref
WHERE env_id = ? AND name = ?;

-- name: VarRefListByVarID :many
SELECT env.name AS env_name, var_ref.* FROM var_ref
JOIN env ON var_ref.env_id = env.env_id
WHERE var_id = ?
ORDER BY var_ref.name ASC;