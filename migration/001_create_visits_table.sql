CREATE TABLE visits (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(100) NOT NULL UNIQUE,
    user_id VARCHAR(100),
    ip INET,
    country VARCHAR(50),
    city VARCHAR(50),
    os VARCHAR(50),
    start_time TIMESTAMP NOT NULL,
    duration FLOAT,
    active_duration FLOAT,
    actions_count INT
);