package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB        *database.DbConfig
	jwtSecret string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env var is not set")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if dbUrl == "" {
		log.Fatal("SERVER_ADDRESS env var is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env var is not set")
	}

	db, err := database.CreateDatabaseConnection(dbUrl)
	if err != nil {
		log.Fatalf("unable to initialize database: %v", err)
	}

	apiCfg := apiConfig{
		DB:        db,
		jwtSecret: jwtSecret,
	}

	router := chi.NewRouter()

	router.Post("/api/toplists", apiCfg.handlerToplistsCreate)
	router.Put("/api/toplists", apiCfg.handlerToplistsUpdate)
	router.Get("/api/toplists/{toplistID}", apiCfg.handlerToplistsGetOne)
	router.Get("/api/toplists", apiCfg.handlerToplistsGetMany)
	router.Delete("/api/toplists/{toplistID}", apiCfg.handlerToplistsDelete)

	router.Post("/api/users", apiCfg.handlerUsersCreate)
	router.Put("/api/users", apiCfg.handlerUsersUpdate)

	router.Post("/api/login", apiCfg.handlerLogin)
	router.Post("/api/refresh", apiCfg.handlerRefresh)
	router.Post("/api/revoke", apiCfg.handlerRevoke)

	srv := &http.Server{
		Handler: router,
		Addr:    serverAddress,
	}

	log.Printf("Server listening on %s", serverAddress)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
