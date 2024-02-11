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

// FeedStore represents a store for managing feed data.
type FeedStore interface {
	CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error)
	ListFeeds(ctx context.Context, arg database.ListFeedsParams) ([]database.Feed, error)
}

// FeedHandler is the handler for feed related requests.
type FeedHandler struct {
	store FeedStore
}

// NewFeedHandler returns a new feed handler.
func NewFeedHandler(store FeedStore) *FeedHandler {
	return &FeedHandler{store: store}
}

// createUserReq is the request to create a user.
type createFeedReq struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// CreateFeed creates a new feed.
func (h *FeedHandler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	var req createFeedReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Log(r.Context(), slog.LevelInfo, "decode feed: %v", err)
		respond.WithJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)
	if userIDVal == nil || !ok {
		respond.WithJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	feed, err := h.store.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:   req.Name,
		Url:    req.URL,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})

	if err != nil {
		// Error should be filtered here.
		slog.Log(r.Context(), slog.LevelError, "create feed: %v", err)
		respond.WithJSONError(w, http.StatusInternalServerError, err.Error())
	}

	respond.WithJSON(w, http.StatusOK, feed)
}

func (h *FeedHandler) ListFeeds(w http.ResponseWriter, r *http.Request) {
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

	feeds, err := h.store.ListFeeds(ctx, database.ListFeedsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		// Error should be filtered here.
		slog.Log(r.Context(), slog.LevelError, "list feeds: %v", err)
		respond.WithJSONError(w, http.StatusInternalServerError, err.Error())
	}

	respond.WithJSON(w, http.StatusOK, feeds)
}
