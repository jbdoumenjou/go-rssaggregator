-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL default now(),
    updated_at TIMESTAMPTZ NOT NULL default now()
);

-- +goose Down
DROP TABLE feed_follows;
