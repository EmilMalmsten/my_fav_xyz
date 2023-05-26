package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
)

func areRanksInOrder(toplistItems []database.ToplistItem) bool {
	for i := 0; i < len(toplistItems); i++ {
		if toplistItems[i].Rank != i+1 {
			return false
		}
	}
	return true
}

func (cfg apiConfig) handlerToplistsAddItems(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ToplistItems []database.ToplistItem `json:"toplistItems"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode json data")
		return
	}

	listIdString := chi.URLParam(r, "listId")
	listId, err := strconv.Atoi(listIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	if !areRanksInOrder(params.ToplistItems) {
		respondWithError(w, http.StatusBadRequest, "Item ranks are not in order")
		return
	}

	addedToplistItems, err := cfg.DB.AddItemsToToplist(params.ToplistItems, listId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error - could not add new items")
		return
	}

	respondWithJSON(w, http.StatusOK, addedToplistItems)
}
