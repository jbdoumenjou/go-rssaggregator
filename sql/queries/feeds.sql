-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY updated_at DESC
LIMIT $1
OFFSET $2;