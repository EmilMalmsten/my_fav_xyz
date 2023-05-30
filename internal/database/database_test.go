package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
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

	db, err := Init(testDbUrl)
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

func createToplist(t *testing.T) int64 {
	arg := Toplist{
		ID:          1,
		Title:       "My Toplist",
		Description: "This is a mock toplist",
		Items: []ToplistItem{
			{
				ID:          101,
				ListId:      1,
				Rank:        1,
				Title:       "Item 1",
				Description: "Description 1",
			},
			{
				ID:          102,
				ListId:      1,
				Rank:        2,
				Title:       "Item 2",
				Description: "Description 2",
			},
		},
	}

	id, err := dbTestConfig.CreateToplist(arg)
	require.NoError(t, err)
	require.NotZero(t, id)

	return id
}

func TestCreateToplist(t *testing.T) {
	createToplist(t)
}
