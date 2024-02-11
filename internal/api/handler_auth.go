package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// UserStore represents a store for managing user data.
type UserStore interface {
	GetUserFromApiKey(ctx context.Context, apiKey string) (database.User, error)
}

// AuthHandler represents an HTTP API handler for user authentication.
type AuthHandler struct {
	store UserStore
}

// NewAuthMiddleware creates a new AuthHandler.
func NewAuthMiddleware(store UserStore) *AuthHandler {
	return &AuthHandler{store: store}
}

func (a *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getAuthHeader(r.Header)
		if err != nil || token == "" {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}

		user, err := a.store.GetUserFromApiKey(r.Context(), token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}

			// Error should be filtered here.
			slog.Log(r.Context(), slog.LevelError, "get user: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user.ID)))
	}
}

func getAuthHeader(h http.Header) (string, error) {
	authHeader := h.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	split := strings.Split(authHeader, " ")
	if strings.ToLower(split[0]) != "apikey" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return split[1], nil
}
