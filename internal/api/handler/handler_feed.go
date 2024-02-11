package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	json2 "github.com/jbdoumenjou/go-rssaggregator/internal/api/json"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// FeedStore represents a store for managing feed data.
type FeedStore interface {
	CreateFeed(ctx context.Context, params database.CreateFeedParams) (database.Feed, error)
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
		json2.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)
	if userIDVal == nil || !ok {
		json2.RespondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
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
		json2.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	json2.RespondWithJSON(w, http.StatusOK, feed)
}
