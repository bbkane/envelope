-- name: VarCreate :exec
INSERT INTO var(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: VarDelete :exec
DELETE FROM var WHERE env_id = ? AND name = ?;

-- name: VarFindID :one
SELECT var_id FROM var WHERE env_id = ? AND name = ?;

-- name: VarFindByID :one
SELECT env.name AS env_name, var.*
FROM var
JOIN env ON var.env_id = env.env_id
WHERE var.var_id = ?;

-- name: VarList :many
SELECT * FROM var
WHERE env_id = ?
ORDER BY name ASC;

-- name: VarShow :one
SELECT *
FROM var
WHERE env_id = ? AND name = ?;

-- name: VarUpdate :exec
UPDATE var SET
    env_id = COALESCE(sqlc.narg('env_id'), env_id),
    name = COALESCE(sqlc.narg('name'), name),
    comment = COALESCE(sqlc.narg('comment'), comment),
    create_time = COALESCE(sqlc.narg('create_time'), create_time),
    update_time = COALESCE(sqlc.narg('update_time'), update_time),
    value = COALESCE(sqlc.narg('value'), value)
WHERE var_id = sqlc.arg('var_id');