package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// FeedFollowsStore represents a store for managing feed follows data.
type FeedFollowsStore interface {
	CreateFeedFollows(ctx context.Context, arg database.CreateFeedFollowsParams) (database.FeedFollow, error)
}

// FeedFollowsHandler is the handler for feed follows related requests.
type FeedFollowsHandler struct {
	store FeedFollowsStore
}

// NewFeedFollowsHandler returns a new feed follows handler.
func NewFeedFollowsHandler(store FeedFollowsStore) *FeedFollowsHandler {
	return &FeedFollowsHandler{store: store}
}

// createUserReq is the request to create a user.
type createFeedFollowsReq struct {
	FeedID uuid.UUID `json:"feed_id"`
}

// CreateFeedFollows creates a new feed follows.
func (h *FeedFollowsHandler) CreateFeedFollows(w http.ResponseWriter, r *http.Request) {
	var req createFeedFollowsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Log(r.Context(), slog.LevelInfo, "decode feed: %v", err)
		respond.WithJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := GetUserIDFromContext(w, r)
	if err != nil {
		slog.Log(r.Context(), slog.LevelInfo, "get user id from context: %v", err)
		respond.WithJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	feedFollows, err := h.store.CreateFeedFollows(r.Context(),
		database.CreateFeedFollowsParams{
			UserID: uuid.NullUUID{
				UUID:  userID,
				Valid: true,
			},
			FeedID: uuid.NullUUID{
				UUID:  req.FeedID,
				Valid: true,
			},
		})
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "create feed follow: %v", err)
		respond.WithJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respond.WithJSON(w, http.StatusOK, feedFollows)
}
