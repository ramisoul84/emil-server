CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    time TIMESTAMP WITH TIME ZONE,
    unread BOOLEAN NOT NULL DEFAULT TRUE,
    country VARCHAR(100)
);

-- Indexes for better performance
CREATE INDEX idx_messages_unread ON messages(unread) WHERE unread = TRUE;