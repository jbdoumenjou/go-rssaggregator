-- name: CreateUser :one
INSERT INTO users (name)
VALUES ($1)
RETURNING *;

-- name: GetUserFromApiKey :one
SELECT * FROM users WHERE api_key = $1;
