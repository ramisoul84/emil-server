CREATE TABLE IF NOT EXISTS visitors (
    id UUID PRIMARY KEY,
    ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    time TIMESTAMP WITH TIME ZONE
);
