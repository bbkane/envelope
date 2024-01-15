INSERT INTO env_var_local(
    env_id, name, comment, create_time, update_time, value
)
SELECT env_id, name, comment, create_time, update_time, local_value as value
FROM env_var;

DROP TABLE env_var;