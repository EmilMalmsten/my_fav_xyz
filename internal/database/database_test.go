package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var dbTestConfig *DbConfig

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function (run)
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {

	godotenv.Load("../../.env")
	testDbUrl := os.Getenv("TEST_DB_URL")
	if testDbUrl == "" {
		log.Fatal("TEST_DB_URL env var is not set")
	}

	db, err := CreateDatabaseConnection(testDbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbTestConfig = db

	tables := []string{"toplists", "list_items"}

	defer func() {
		for _, t := range tables {
			_, _ = dbTestConfig.database.Exec(fmt.Sprintf("DELETE FROM %s", t))
		}
		dbTestConfig.database.Close()
	}()

	return m.Run(), nil
}
