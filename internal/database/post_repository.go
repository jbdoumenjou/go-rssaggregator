package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// PostRepository is responsible for managing the users in the database.
type PostRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewPostRepository creates a new PostRepository.
func NewPostRepository(db *sql.DB) PostRepository {
	return PostRepository{
		db:      db,
		queries: New(db),
	}
}

// GetPostsByUser returns the posts for the given user.
func (u PostRepository) GetPostsByUser(ctx context.Context, userID uuid.NullUUID, limit int32) ([]Post, error) {
	posts, err := u.queries.GetPostsByUser(ctx, GetPostsByUserParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting posts by user %w", err)
	}

	return posts, nil
}

// CreatePost creates a new post.
func (u PostRepository) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	post, err := u.queries.CreatePost(ctx, arg)
	if err != nil {
		return Post{}, fmt.Errorf("error creating post: %w", err)
	}

	return post, nil
}
