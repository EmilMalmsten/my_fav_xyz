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
	dbUrl string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env var is not set")
	}

	apiCfg := apiConfig{
		dbUrl: dbUrl,
	}

	router := chi.NewRouter()

	router.Post("/api/toplists", apiCfg.handlerToplistsCreate)

	database.Db(apiCfg.dbUrl)

	fmt.Println("Server running...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
