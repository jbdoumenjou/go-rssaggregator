package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	// don't forget to import the driver ;)
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain is the entry point for the test suite.
// It is called only once, before all tests.
// It is responsible for initializing the test database.
func TestMain(m *testing.M) {
	var err error
	fmt.Println("Initializing test suite")

	if err := godotenv.Load(); err != nil {
		log.New(os.Stderr, "Error loading .env file", 0)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL env variable not set")
	}

	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to the db:", err)
	}
	testQueries = New(testDB)
	_, err = testDB.Exec("TRUNCATE TABLE users")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
