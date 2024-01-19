-- name: CreateKeyringEntry :exec
INSERT INTO keyring_entry(
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
);

-- name: FindKeyringID :one
SELECT id from keyring_entry WHERE name = ?