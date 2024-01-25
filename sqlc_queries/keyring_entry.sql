-- name: KeyringEntryCreate :exec
INSERT INTO keyring_entry(
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
);

-- name: KeyringEntryFindID :one
SELECT id from keyring_entry WHERE name = ?;

-- name: KeyringEntryShow :one
SELECT * FROM keyring_entry WHERE name = ?;