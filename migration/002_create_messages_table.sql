CREATE TABLE messages(
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    text TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    unread BOOLEAN NOT NULL,
    ip INET,
    city VARCHAR(50),
    country VARCHAR(50)
);