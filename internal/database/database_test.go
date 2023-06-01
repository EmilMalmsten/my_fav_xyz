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

func insertToplist(t *testing.T) Toplist {
	toplist := Toplist{
		Title:       "My Toplist",
		Description: "This is a mock toplist",
		Items: []ToplistItem{
			{
				Rank:        1,
				Title:       "Item 1",
				Description: "Description 1",
			},
			{
				Rank:        2,
				Title:       "Item 2",
				Description: "Description 2",
			},
		},
	}

	insertedToplist, err := dbTestConfig.InsertToplist(toplist)
	require.NoError(t, err)
	require.NotZero(t, insertedToplist)

	require.Equal(t, toplist.Title, insertedToplist.Title)
	require.Equal(t, toplist.Description, insertedToplist.Description)

	require.Equal(t, len(toplist.Items), len(insertedToplist.Items))

	for i := range insertedToplist.Items {
		require.Equal(t, insertedToplist.Items[i].ListID, insertedToplist.ID)
		require.Equal(t, toplist.Items[i].Rank, insertedToplist.Items[i].Rank)
		require.Equal(t, toplist.Items[i].Title, insertedToplist.Items[i].Title)
		require.Equal(t, toplist.Items[i].Description, insertedToplist.Items[i].Description)
	}

	return insertedToplist
}

func TestInsertToplist(t *testing.T) {
	insertToplist(t)
}

func TestGetToplist(t *testing.T) {
	toplist1 := insertToplist(t)
	toplist2, err := dbTestConfig.GetToplist(toplist1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, toplist2)

	require.Equal(t, toplist1.ID, toplist2.ID)
	require.Equal(t, toplist1.Title, toplist2.Title)
	require.Equal(t, toplist1.Description, toplist2.Description)

	require.Equal(t, len(toplist1.Items), len(toplist2.Items))

	for i := range toplist2.Items {
		require.Equal(t, toplist1.Items[i].ListID, toplist2.Items[i].ListID)
		require.Equal(t, toplist1.Items[i].Rank, toplist2.Items[i].Rank)
		require.Equal(t, toplist1.Items[i].Title, toplist2.Items[i].Title)
		require.Equal(t, toplist1.Items[i].Description, toplist2.Items[i].Description)
	}
}

func TestUpdateToplist(t *testing.T) {
	toplist1 := insertToplist(t)

	toplist2 := toplist1
	toplist2.Title = "Updated My Toplist"

	toplist2, err := dbTestConfig.UpdateToplist(toplist2)
	require.NoError(t, err)
	require.NotEmpty(t, toplist2)

	require.NotEqual(t, toplist1.Title, toplist2.Title)
	require.Equal(t, toplist1.Description, toplist2.Description)
	require.Equal(t, toplist1.ID, toplist2.ID)
}

func TestRemoveToplist(t *testing.T) {
	toplist := insertToplist(t)

	err := dbTestConfig.DeleteToplist(toplist.ID)
	require.NoError(t, err)

	_, err = dbTestConfig.GetToplist(toplist.ID)
	require.ErrorIs(t, err, ErrNotExist)
}
