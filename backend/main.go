package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)

	router.With(apiCfg.validateJWT).Post("/api/toplists", apiCfg.handlerToplistsCreate)
	router.With(apiCfg.validateJWT).Put("/api/toplists", apiCfg.handlerToplistsUpdate)
	router.Get("/api/toplists/{toplistID}", apiCfg.handlerToplistsGetOne)
	router.Get("/api/toplists", apiCfg.handlerToplistsGetMany)
	router.Get("/api/toplists/recent", apiCfg.handlerToplistsGetRecent)
	router.With(apiCfg.validateJWT).Delete("/api/toplists/{toplistID}", apiCfg.handlerToplistsDelete)

	router.Post("/api/users", apiCfg.handlerUsersCreate)
	router.With(apiCfg.validateJWT).Put("/api/users", apiCfg.handlerUsersUpdate)
	router.With(apiCfg.validateJWT).Delete("/api/users", apiCfg.handlerUsersDelete)

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
