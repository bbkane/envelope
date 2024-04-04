-- name: KeyringEntryCreate :exec
INSERT INTO keyring_entry(
    name, comment, create_time, update_time
) VALUES (
    ?   , ?      , ?          , ?
);

-- name: KeyringEntryDelete :exec
DELETE FROM keyring_entry WHERE name = ?;

-- name: KeyringEntryList :many
SELECT * FROM keyring_entry
ORDER BY name ASC;

-- name: KeyringEntryFindID :one
SELECT keyring_entry_id from keyring_entry WHERE name = ?;

-- name: KeyringEntryShow :one
SELECT * FROM keyring_entry WHERE name = ?;