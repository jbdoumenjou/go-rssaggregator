package database

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueries_CreateFeedFollows(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	feed := CreateRandomFeed(t)

	follows, err := testQueries.CreateFeedFollows(ctx, CreateFeedFollowsParams{
		UserID: feed.UserID,
		FeedID: uuid.NullUUID{UUID: feed.ID, Valid: true},
	})
	require.NoError(t, err)
	require.NotEmpty(t, follows)
	assert.Equal(t, feed.UserID, follows.UserID)
	assert.Equal(t, feed.ID, follows.FeedID.UUID)
	assert.NotEmpty(t, follows.CreatedAt)
	assert.NotEmpty(t, follows.UpdatedAt)
}

func TestQueries_CreateFeedFollows_UnknownUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	feed := CreateRandomFeed(t)

	follows, err := testQueries.CreateFeedFollows(ctx, CreateFeedFollowsParams{
		UserID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
		FeedID: uuid.NullUUID{UUID: feed.ID, Valid: true},
	})
	require.Error(t, err)
	require.Empty(t, follows)
}

func TestQueries_CreateFeedFollows_UnknownFeed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	feed := CreateRandomFeed(t)

	follows, err := testQueries.CreateFeedFollows(ctx, CreateFeedFollowsParams{
		UserID: feed.UserID,
		FeedID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
	})
	require.Error(t, err)
	require.Empty(t, follows)
}
