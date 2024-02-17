package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// FeedFollowsStore represents a store for managing feed follows data.
type FeedFollowsStore interface {
	CreateFeedFollows(ctx context.Context, arg database.CreateFeedFollowsParams) (database.FeedFollow, error)
	ListFeedFollows(ctx context.Context, arg database.ListFeedFollowsParams) ([]database.FeedFollow, error)
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

// ListFeedFollows lists the feed follows.
func (h *FeedFollowsHandler) ListFeedFollows(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)
	if userIDVal == nil || !ok {
		respond.WithJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	// Get the values of 'offset' and 'limit' from the URL query parameters
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		respond.WithJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid offset: %q", offsetStr))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		respond.WithJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid limit: %q", limitStr))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	feeds, err := h.store.ListFeedFollows(ctx, database.ListFeedFollowsParams{
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		// Error should be filtered here.
		slog.Log(r.Context(), slog.LevelError, "list feeds: %v", err)
		respond.WithJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.WithJSON(w, http.StatusOK, feeds)
}
