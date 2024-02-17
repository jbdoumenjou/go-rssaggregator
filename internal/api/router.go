package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/handler"
	"github.com/jbdoumenjou/go-rssaggregator/internal/api/middleware"
)

type Router struct {
	mux *chi.Mux

	authHandler        *middleware.AuthHandler
	userHandler        *handler.UserHandler
	feedHandler        *handler.FeedHandler
	feedFollowsHandler *handler.FeedFollowsHandler
}

func NewRouter(authHandler *middleware.AuthHandler, userHandler *handler.UserHandler, feedHandler *handler.FeedHandler, feedFollowsHandler *handler.FeedFollowsHandler) http.Handler {
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

	router := &Router{
		mux:                r,
		authHandler:        authHandler,
		userHandler:        userHandler,
		feedHandler:        feedHandler,
		feedFollowsHandler: feedFollowsHandler,
	}

	router.addV1Routes()

	return r
}

func (r Router) addV1Routes() {
	v1 := chi.NewRouter()
	r.mux.Mount("/v1", v1)

	v1.Get("/readiness", handler.Readiness)
	v1.Get("/err", handler.Error)

	v1.Post("/users", r.userHandler.CreateUser)
	v1.Get("/users", r.authHandler.Authenticate(r.userHandler.GetUser))

	v1.Post("/feeds", r.authHandler.Authenticate(r.feedHandler.CreateFeed))
	v1.Get("/feeds", r.feedHandler.ListFeeds)

	v1.Post("/feed_follows", r.authHandler.Authenticate(r.feedFollowsHandler.CreateFeedFollows))
	v1.Get("/feed_follows", r.authHandler.Authenticate(r.feedFollowsHandler.ListFeedFollows))
	v1.Delete("/feed_follows/{id}", r.authHandler.Authenticate(r.feedFollowsHandler.DeleteFeedFollows))
}
