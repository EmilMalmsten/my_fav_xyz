package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)

	router.With(apiCfg.validateJWT).Post("/api/toplists", apiCfg.handlerToplistsCreate)
	router.With(apiCfg.validateJWT).Put("/api/toplists/items", apiCfg.handlerToplistsUpdateItems)
	router.With(apiCfg.validateJWT).Put("/api/toplists", apiCfg.handlerToplistsUpdate)
	router.With(apiCfg.validateJWT).Delete("/api/toplists/{toplistID}", apiCfg.handlerToplistsDelete)
	router.Get("/api/toplists/{toplistID}", apiCfg.handlerToplistsGetOne)
	router.Get("/api/toplists", apiCfg.handlerToplistsGetMany)
	router.Get("/api/toplists/recent", apiCfg.handlerToplistsGetRecent)
	router.Get("/api/toplists/popular", apiCfg.handlerToplistsGetPopular)
	router.Post("/api/toplists/views/{toplistID}", apiCfg.handlerToplistsViews)

	router.Post("/api/users", apiCfg.handlerUsersCreate)
	router.Put("/api/users/email", apiCfg.handlerUsersUpdateEmail)
	router.Put("/api/users/password", apiCfg.handlerUsersUpdatePassword)
	router.Delete("/api/users", apiCfg.handlerUsersDelete)

	router.Post("/api/login", apiCfg.handlerLogin)
	router.Post("/api/refresh", apiCfg.handlerRefresh)
	router.Post("/api/revoke", apiCfg.handlerRevoke)

	FileServer(router, "/images", http.Dir("./internal/database/images"))

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

func FileServer(r chi.Router, path string, root http.FileSystem) {
    if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}