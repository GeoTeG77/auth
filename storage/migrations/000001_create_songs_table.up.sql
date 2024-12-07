CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    refresh_token_hash TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_refresh_token_hash ON users USING hash (refresh_token_hash);