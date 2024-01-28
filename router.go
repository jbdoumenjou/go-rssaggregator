package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

func NewRouter(db database.Querier) http.Handler {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Adds a v1 subrouter.
	v1 := chi.NewRouter()
	r.Mount("/v1", v1)
	addV1Routes(v1, db)

	return r
}

func addV1Routes(r chi.Router, db database.Querier) {
	r.Get("/readiness", readinessHandler)
	r.Get("/err", errorHandler)

	userHandler := NewUserHandler(db)
	r.Post("/users", userHandler.CreateUser)
}
