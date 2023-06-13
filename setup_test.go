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

	r := chi.NewRouter()
	r.Put("/api/toplists", apiCfg.handlerToplistsCreate)

	code := m.Run()

	os.Exit(code)
}
