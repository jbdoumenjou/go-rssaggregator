package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// UserRepository is responsible for managing the users in the database.
type UserRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db:      db,
		queries: New(db),
	}
}

// GetUserFromApiKey returns the user with the given api key.
func (u UserRepository) GetUserFromApiKey(ctx context.Context, apiKey string) (User, error) {
	user, err := u.queries.GetUserFromApiKey(ctx, apiKey)
	if err != nil {
		return User{}, fmt.Errorf("error getting user from api key: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user.
func (u UserRepository) CreateUser(ctx context.Context, name string) (User, error) {
	user, err := u.queries.CreateUser(ctx, name)
	if err != nil {
		return User{}, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// GetUserFromId returns the user with the given id.
func (u UserRepository) GetUserFromId(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := u.queries.GetUserFromId(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("error getting user from id: %w", err)
	}

	return user, nil
}
