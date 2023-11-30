package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

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

	// Create a new router.
	r := NewRouter()

	// start the server.
	if err := NewServer("localhost:"+port, r).Start(); err != nil {
		log.Fatal(err)
	}
}
