package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
)

func (cfg apiConfig) handlerToplistsChangeItems(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ToplistItems []database.ToplistItem `json:"toplistItems"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode json data")
		return
	}
	fmt.Println(params.ToplistItems)

	listIDString := chi.URLParam(r, "listId")
	listID, err := strconv.Atoi(listIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	if !areRanksInOrder(params.ToplistItems) {
		respondWithError(w, http.StatusBadRequest, "Item ranks are not in order")
		return
	}

	err = cfg.DB.AddItemsToToplist(params.ToplistItems, listID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add items to list")
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}

func areRanksInOrder(toplistItems []database.ToplistItem) bool {
	for i := 0; i < len(toplistItems); i++ {
		if toplistItems[i].Rank != i+1 {
			return false
		}
	}
	return true
}
