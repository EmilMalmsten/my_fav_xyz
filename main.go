package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type apiConfig struct {
}

func main() {
	apiCfg := apiConfig{}

	router := chi.NewRouter()

	router.Post("/api/toplists", apiCfg.handlerCreateToplist)

	fmt.Println("Server running...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
