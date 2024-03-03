// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	Url           string        `json:"url"`
	UserID        uuid.NullUUID `json:"user_id"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	LastFetchedAt sql.NullTime  `json:"last_fetched_at"`
}

type FeedFollow struct {
	ID        uuid.UUID     `json:"id"`
	FeedID    uuid.NullUUID `json:"feed_id"`
	UserID    uuid.NullUUID `json:"user_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Post struct {
	ID          int32         `json:"id"`
	Title       string        `json:"title"`
	Url         string        `json:"url"`
	Description string        `json:"description"`
	PublishedAt time.Time     `json:"published_at"`
	FeedID      uuid.NullUUID `json:"feed_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
}
