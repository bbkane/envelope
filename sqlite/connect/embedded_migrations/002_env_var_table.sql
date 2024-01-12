CREATE TABLE env_var (
    id INTEGER PRIMARY KEY,
    env_id INTEGER NOT NULL REFERENCES env(id),
    name TEXT NOT NULL UNIQUE,
    comment TEXT,
    create_time TEXT NOT NULL,
    update_time TEXT NOT NULL,
    type TEXT NOT NULL,
    local_value TEXT,
    FOREIGN KEY (env_id) REFERENCES env(id),
    UNIQUE(env_id, name)
) STRICT;