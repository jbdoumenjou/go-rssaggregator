package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// FeedRepository is responsible for managing the feeds in the database.
type FeedRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewFeedRepository creates a new FeedRepository.
func NewFeedRepository(db *sql.DB) FeedRepository {
	return FeedRepository{
		db:      db,
		queries: New(db),
	}
}

// CreateFeed creates a new feed.
func (f FeedRepository) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	feed, err := f.queries.CreateFeed(ctx, arg)
	if err != nil {
		return Feed{}, fmt.Errorf("error creating feed: %w", err)
	}

	return feed, nil
}

// CreateFeedAndFollow creates a new feed and follows it.
func (f FeedRepository) CreateFeedAndFollow(ctx context.Context, arg CreateFeedParams) (Feed, FeedFollow, error) {
	tx, err := f.db.Begin()
	if err != nil {
		return Feed{}, FeedFollow{}, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := f.queries.WithTx(tx)
	feed, err := qtx.CreateFeed(ctx, arg)
	if err != nil {
		return Feed{}, FeedFollow{}, fmt.Errorf("error creating feed: %w", err)
	}
	follow, err := qtx.CreateFeedFollows(ctx, CreateFeedFollowsParams{
		UserID: arg.UserID,
		FeedID: uuid.NullUUID{
			UUID:  feed.ID,
			Valid: true,
		},
	})
	if err != nil {
		return Feed{}, FeedFollow{}, fmt.Errorf("error creating feed follow: %w", err)
	}

	return feed, follow, tx.Commit()
}

// ListFeeds returns a list of feeds.
func (f FeedRepository) ListFeeds(ctx context.Context, arg ListFeedsParams) ([]Feed, error) {
	feeds, err := f.queries.ListFeeds(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("error listing feeds: %w", err)
	}

	return feeds, nil
}

// CreateFeedFollows creates a new feed follow.
func (f FeedRepository) CreateFeedFollows(ctx context.Context, arg CreateFeedFollowsParams) (FeedFollow, error) {
	follow, err := f.queries.CreateFeedFollows(ctx, arg)
	if err != nil {
		return FeedFollow{}, fmt.Errorf("error creating feed follow: %w", err)
	}

	return follow, nil
}

// ListFeedFollows returns a list of feed follows.
func (f FeedRepository) ListFeedFollows(ctx context.Context, arg ListFeedFollowsParams) ([]FeedFollow, error) {
	follows, err := f.queries.ListFeedFollows(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("error listing feed follows: %w", err)
	}

	return follows, nil
}

// DeleteFeedFollows deletes a feed follow.
func (f FeedRepository) DeleteFeedFollows(ctx context.Context, arg DeleteFeedFollowsParams) error {
	err := f.queries.DeleteFeedFollows(ctx, arg)
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %w", err)
	}

	return nil
}
