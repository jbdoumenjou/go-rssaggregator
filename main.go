package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api/middleware"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api/handler"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load .env file variables.
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("PORT env variable not set, using default port %s\n", "8080")
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL env variable not set")
	}

	// Connect to the database.
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to the db:", err)
	}
	userRepository := database.NewUserRepository(db)
	feedRepository := database.NewFeedRepository(db)

	authHandler := middleware.NewAuthMiddleware(userRepository)
	userHandler := handler.NewUserHandler(userRepository)
	feedHandler := handler.NewFeedHandler(feedRepository)
	feedFollowsHandler := handler.NewFeedFollowsHandler(feedRepository)

	// Create a new router.
	r := api.NewRouter(authHandler, userHandler, feedHandler, feedFollowsHandler)

	// start the server.
	if err := api.NewServer("localhost:"+port, r).Start(); err != nil {
		log.Fatal(err)
	}
}
