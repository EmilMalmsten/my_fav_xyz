package main

import (
	"encoding/json"
	"net/http"
)

func (cfg apiConfig) handlerToplistsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	type resp struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	toplistId, err := cfg.DB.CreateToplist(params.Title, params.Description)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new toplist")
		return
	}

	respondWithJSON(w, http.StatusCreated, resp{
		Id:          toplistId,
		Title:       params.Title,
		Description: params.Description,
	})
}
