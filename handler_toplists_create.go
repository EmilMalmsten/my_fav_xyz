package main

import (
	"encoding/json"
	"net/http"
)

func (cfg apiConfig) handlerToplistsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	respondWithJSON(w, http.StatusCreated, params)
}
