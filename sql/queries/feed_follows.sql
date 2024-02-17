-- name: CreateFeedFollows :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows
WHERE id = $1;

-- name: ListFeedFollows :many
SELECT * FROM feed_follows
WHERE user_id = $1
ORDER BY updated_at DESC
LIMIT $2
OFFSET $3;