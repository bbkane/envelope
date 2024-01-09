-- name: CreateEnv :one
INSERT INTO env (
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
)
RETURNING id;

-- name: UpdateEnv :exec
UPDATE env SET
    name = COALESCE(?, name),
    comment = COALESCE(?, comment),
    create_time = COALESCE(?, create_time),
    update_time = COALESCE(?, update_time)
WHERE id = ?;