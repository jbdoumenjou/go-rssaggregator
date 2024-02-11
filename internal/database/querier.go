// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error)
	CreateFeedFollows(ctx context.Context, arg CreateFeedFollowsParams) (FeedFollow, error)
	CreateUser(ctx context.Context, name string) (User, error)
	DeleteFeedFollows(ctx context.Context, id uuid.UUID) error
	GetFeedFollows(ctx context.Context, userID uuid.NullUUID) ([]FeedFollow, error)
	GetUserFromApiKey(ctx context.Context, apiKey string) (User, error)
	GetUserFromId(ctx context.Context, id uuid.UUID) (User, error)
	ListFeeds(ctx context.Context, arg ListFeedsParams) ([]Feed, error)
}

var _ Querier = (*Queries)(nil)
