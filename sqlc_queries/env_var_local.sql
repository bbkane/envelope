-- name: CreateLocalEnvVar :exec
INSERT INTO env_var_local(
    env_id, name, comment, create_time, update_time, value
) VALUES (
    ?     , ?   , ?      , ?          , ?          , ?
);

-- name: ListLocalEnvVars :many
SELECT * FROM env_var_local
WHERE env_id = ?
ORDER BY name ASC;