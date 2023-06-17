package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

var apiCfg apiConfig
var insertedTestToplists []database.Toplist

func TestMain(m *testing.M) {
	fmt.Println("Initializing test...")
	godotenv.Load()
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		log.Fatal("TEST_DB_URL env var is not set")
	}

	db, err := database.CreateDatabaseConnection(dbURL)
	if err != nil {
		log.Fatalf("unable to initialize database: %v", err)
	}

	apiCfg = apiConfig{
		DB: db,
	}

	insertToplists()

	r := chi.NewRouter()
	r.Put("/api/toplists", apiCfg.handlerToplistsCreate)
	r.Delete("/api/toplists/{toplistID}", apiCfg.handlerToplistsDelete)

	code := m.Run()

	os.Exit(code)
}

func insertToplists() {
	toplists := []toplistRequest{
		{
			Title:       "Test toplist 1",
			Description: "Test description 1",
			UserID:      1,
			Items: []toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
		},
		{
			Title:       "Test toplist 2",
			Description: "Test description 2",
			UserID:      1,
			Items: []toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
		},
		{
			Title:       "Test toplist 3",
			Description: "Test description 3",
			UserID:      1,
			Items: []toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
		},
	}

	for _, toplist := range toplists {
		dbToplist := toplist.ToDBToplist()
		insertedToplist, err := apiCfg.DB.InsertToplist(dbToplist)
		if err != nil {
			log.Fatalf("unable to insert test toplist: %v", err)
		}
		insertedTestToplists = append(insertedTestToplists, insertedToplist)
	}
}
