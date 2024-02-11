package database

import (
	"context"
	"testing"

	"github.com/jbdoumenjou/go-rssaggregator/internal/generator"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	t.Helper()

	name := generator.RandomString(12)
	user, err := testQueries.CreateUser(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEmpty(t, user.ID)
	require.Equal(t, name, user.Name)
	require.NotEmpty(t, user.CreatedAt)
	require.NotEmpty(t, user.UpdatedAt)
	require.NotEmpty(t, user.ApiKey)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	user := CreateRandomUser(t)

	testQueries.db.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)
}

func TestQueries_GetUserFromApiKey(t *testing.T) {
	user := CreateRandomUser(t)

	user2, err := testQueries.GetUserFromApiKey(context.Background(), user.ApiKey)
	require.NoError(t, err)
	assert.Equal(t, user, user2)

	testQueries.db.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)

	user3, err := testQueries.GetUserFromApiKey(context.Background(), user.ApiKey)
	require.Error(t, err)
	assert.Empty(t, user3)
}

func TestQueries_DeleteUserDeleteFeeds(t *testing.T) {
	feed := CreateRandomFeed(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query1, err := testQueries.db.QueryContext(ctx, "Delete from users where id = $1", feed.UserID.UUID)
	require.NoError(t, err)
	defer query1.Close()

	query2, err := testQueries.db.QueryContext(ctx, "SELECT from feeds where id = $1", feed.ID)
	require.NoError(t, err)
	defer query2.Close()
	// feeds from the user should be deleted
	assert.False(t, query2.Next())
}

func TestQueries_GetUserFromId(t *testing.T) {
	user := CreateRandomUser(t)

	user2, err := testQueries.GetUserFromId(context.Background(), user.ID)
	require.NoError(t, err)
	assert.Equal(t, user, user2)

	testQueries.db.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)

	user3, err := testQueries.GetUserFromId(context.Background(), user.ID)
	require.Error(t, err)
	assert.Empty(t, user3)
}
