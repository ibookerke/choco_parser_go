CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);