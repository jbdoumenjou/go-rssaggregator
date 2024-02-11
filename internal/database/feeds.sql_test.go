package database

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomFeed(t *testing.T) Feed {
	t.Helper()

	user := CreateRandomUser(t)
	require.NotEmpty(t, user)
	require.NotEmpty(t, user.ID)

	feedParams := CreateFeedParams{
		Name:   RandomString(12),
		Url:    fmt.Sprintf("https://%s.%s", RandomString(8), RandomString(3)),
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed, err := testQueries.CreateFeed(ctx, feedParams)
	require.NoError(t, err)
	require.NotEmpty(t, feed)
	assert.Equal(t, feedParams.Name, feed.Name)
	assert.Equal(t, feedParams.Url, feed.Url)
	assert.Equal(t, feedParams.UserID, feed.UserID)
	require.NotEmpty(t, feed.CreatedAt)
	require.NotEmpty(t, feed.UpdatedAt)

	return feed
}

func TestQueries_CreateFeed(t *testing.T) {
	feed := CreateRandomFeed(t)
	assert.NotEmpty(t, feed)
}

func TestQueries_CreateFeed_BadUserID(t *testing.T) {
	feedParams := CreateFeedParams{
		Name:   RandomString(12),
		Url:    fmt.Sprintf("https://%s.%s", RandomString(8), RandomString(3)),
		UserID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed, err := testQueries.CreateFeed(ctx, feedParams)
	assert.Error(t, err)
	assert.Empty(t, feed)
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	require.True(t, ok)
	assert.Equal(t, "foreign_key_violation", pqErr.Code.Name())
}

func TestQueries_CreateFeed_ExistingURL(t *testing.T) {
	feed1 := CreateRandomFeed(t)
	require.NotEmpty(t, feed1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed2, err := testQueries.CreateFeed(ctx, CreateFeedParams{
		Name:   "name2",
		Url:    feed1.Url,
		UserID: feed1.UserID,
	})

	assert.Error(t, err)
	assert.Empty(t, feed2)
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	require.True(t, ok)
	assert.Equal(t, "unique_violation", pqErr.Code.Name())
}
