CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    units TEXT,
    threshold DOUBLE PRECISION,
    is_log BOOLEAN,
    description TEXT,
    cas_type CHAR(2),
    log_type TEXT
);
