package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now().UTC().Round(time.Microsecond)
	user, err := h.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		// Error should be filtered here.
		slog.Log(r.Context(), slog.LevelError, "create user: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, user)
}
