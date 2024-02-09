package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
	"github.com/joho/godotenv"
)

var apiCfg apiConfig
var insertedTestUser database.User
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

	err = insertTestUser()
	if err != nil {
		log.Fatalf("unable to insert test user: %v", err)
	}

	insertTestToplists(insertedTestUser.ID)

	code := m.Run()

	os.Exit(code)
}

func createRandomEmail() string {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	randomNumber := rng.Intn(100000) + 1
	email := fmt.Sprintf("testuser%d@mail.com", randomNumber)
	return email
}

func createRandomUsername() string {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	randomNumber := rng.Intn(100000) + 1
	username := fmt.Sprintf("testuser%d", randomNumber)
	return username
}

func insertTestUser() error {
	randomEmail := createRandomEmail()
	randomUsername := createRandomUsername()
	user := database.User{
		Email:          randomEmail,
		Username:       randomUsername,
		HashedPassword: "asd123123123hjerwehr",
	}
	insertedUser, err := apiCfg.DB.InsertUser(user)
	if err != nil {
		return err
	}

	insertedTestUser = insertedUser
	return nil
}

func insertTestToplists(testUserID int) {
	toplists := []toplistRequest{
		{
			Title:       "Test toplist 1",
			Description: "Test description 1",
			UserID:      testUserID,
			Items: []toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
		},
		{
			Title:       "Test toplist 2",
			Description: "Test description 2",
			UserID:      testUserID,
			Items: []toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
		},
		{
			Title:       "Test toplist 3",
			Description: "Test description 3",
			UserID:      testUserID,
			Items: []toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
		},
		{
			Title:       "Test toplist 4",
			Description: "Test description 4",
			UserID:      testUserID,
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
