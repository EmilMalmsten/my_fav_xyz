package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.DbConfig
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env var is not set")
	}

	db, err := database.CreateDatabaseConnection(dbUrl)
	if err != nil {
		log.Fatalf("unable to initialize database: %v", err)
	}

	apiCfg := apiConfig{
		DB: db,
	}

	router := chi.NewRouter()

	router.Post("/api/toplists", apiCfg.handlerToplistsCreate)
	router.Put("/api/toplists/{listID}", apiCfg.handlerToplistsUpdate)
	router.Put("/api/toplists/{listID}/items", apiCfg.handlerToplistsChangeItems)

	fmt.Println("Server running...")
	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		panic(err)
	}
}
