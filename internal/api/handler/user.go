package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// UserHandler is the handler for user related requests.
type UserHandler struct {
	db database.Querier
}

// NewUserHandler returns a new user handler.
func NewUserHandler(db database.Querier) *UserHandler {
	return &UserHandler{db: db}
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

	user, err := h.db.CreateUser(r.Context(), req.Name)
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

	user, err := h.db.GetUserFromId(r.Context(), userID)
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

func GetUserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	userIDVal := r.Context().Value("user")
	userID, ok := userIDVal.(uuid.UUID)

	if userIDVal == nil || !ok {
		return uuid.UUID{}, errors.New("cannot get user id from context")
	}

	return userID, nil
}
