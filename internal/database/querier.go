// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"context"
)

type Querier interface {
	CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error)
	CreateUser(ctx context.Context, name string) (User, error)
	GetUserFromApiKey(ctx context.Context, apiKey string) (User, error)
}

var _ Querier = (*Queries)(nil)
