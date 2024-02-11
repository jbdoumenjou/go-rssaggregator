-- +goose Up
CREATE TABLE feeds (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name VARCHAR NOT NULL,
                       url VARCHAR UNIQUE NOT NULL,
                       user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                       created_at TIMESTAMPTZ NOT NULL default now(),
                       updated_at TIMESTAMPTZ NOT NULL default now()
);

-- +goose Down
DROP TABLE feeds;
