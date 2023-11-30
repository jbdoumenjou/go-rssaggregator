package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestQueries_CreateUser(t *testing.T) {
	now := time.Now().UTC().Round(time.Microsecond)
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)

	params := CreateUserParams{
		ID:        uuid,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      RandomString(12),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.ID, user.ID)
	require.Equal(t, params.Name, user.Name)
	require.Equal(t, params.CreatedAt, user.CreatedAt)
	require.Equal(t, params.UpdatedAt, user.UpdatedAt)
}
