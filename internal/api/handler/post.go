package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// PostStore represents a store for managing post data.
type PostStore interface {
	GetPostsByUser(ctx context.Context, userID uuid.NullUUID, limit int32) ([]database.Post, error)
}

// PostHandler is the handler for feed related requests.
type PostHandler struct {
	store PostStore
}

// NewPostHandler returns a new feed handler.
func NewPostHandler(store PostStore) *PostHandler {
	return &PostHandler{store: store}
}

// GetPostsByUser returns the posts for the given user.
func (h *PostHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)
	if userIDVal == nil || !ok {
		respond.WithJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
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

	posts, err := h.store.GetPostsByUser(ctx, uuid.NullUUID{UUID: userID, Valid: true}, int32(limit))
	if err != nil {
		respond.WithJSONError(w, http.StatusInternalServerError, "error getting posts")
		return
	}

	respond.WithJSON(w, http.StatusOK, posts)
}
