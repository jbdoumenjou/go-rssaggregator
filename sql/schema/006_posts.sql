-- +goose Up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL default now(),
    updated_at TIMESTAMPTZ NOT NULL default now()
);

-- +goose Down
DROP TABLE posts;


