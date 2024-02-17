package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// UserStore represents a store for managing user data.
type UserStore interface {
	CreateUser(ctx context.Context, name string) (database.User, error)
	GetUserFromId(ctx context.Context, id uuid.UUID) (database.User, error)
}

// UserHandler is the handler for user related requests.
type UserHandler struct {
	store UserStore
}

// NewUserHandler returns a new user handler.
func NewUserHandler(store UserStore) *UserHandler {
	return &UserHandler{store: store}
}

// createUserReq is the request to create a user.
type createUserReq struct {
	Name string `json:"name"`
}

// CreateUser creates a new user.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Log(r.Context(), slog.LevelInfo, "decode user: %v", err)
		respond.WithJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.store.CreateUser(r.Context(), req.Name)
	if err != nil {
		// Error should be filtered here.
		slog.Log(r.Context(), slog.LevelError, "create user: %v", err)
		respond.WithJSONError(w, http.StatusInternalServerError, err.Error())
	}

	respond.WithJSON(w, http.StatusOK, user)
}

// GetUser get a user from the authorization header.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)
	if userIDVal == nil || !ok {
		respond.WithJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	user, err := h.store.GetUserFromId(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.WithJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		// Error should be filtered here.
		slog.Log(r.Context(), slog.LevelError, "get user: %v", err)
		respond.WithJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.WithJSON(w, http.StatusOK, user)
}

// GetUserIDFromContext gets the user id from the context.
func GetUserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)

	if userIDVal == nil || !ok {
		return uuid.UUID{}, errors.New("cannot get user id from context")
	}

	return userID, nil
}
