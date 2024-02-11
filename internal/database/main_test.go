package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

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

	pgContainer, err := NewPGContainer()
	if err != nil {
		log.Fatal("cannot create a new container:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal("cannot terminate the container:", err)
		}
	}()

	testDB = pgContainer.DB()
	testQueries = New(testDB)

	_, err = testDB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
