-- +goose Up
CREATE TABLE feeds_fetched ( 
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    last_fetched_at TIMESTAMP, 
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE feeds;

-- the feeds table, but with an added last_fetched column