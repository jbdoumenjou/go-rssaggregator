package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	t.Helper()

	name := RandomString(12)
	user, err := testQueries.CreateUser(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEmpty(t, user.ID)
	require.Equal(t, name, user.Name)
	require.NotEmpty(t, user.CreatedAt)
	require.NotEmpty(t, user.UpdatedAt)
	require.NotEmpty(t, user.Apikey)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	user := CreateRandomUser(t)

	testQueries.db.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)
}

func TestQueries_GetUserFromApiKey(t *testing.T) {
	user := CreateRandomUser(t)

	user2, err := testQueries.GetUserFromApiKey(context.Background(), user.Apikey)
	require.NoError(t, err)
	assert.Equal(t, user, user2)

	testQueries.db.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)

	user3, err := testQueries.GetUserFromApiKey(context.Background(), user.Apikey)
	require.Error(t, err)
	assert.Empty(t, user3)
}
