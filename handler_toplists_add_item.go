package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (cfg apiConfig) handlerToplistsAddItem(w http.ResponseWriter, r *http.Request) {
	type ToplistItem struct {
		Rank        int    `json:"rank"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	type parameters struct {
		ToplistItems []ToplistItem `json:"toplistItems"`
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

	fmt.Println(listId)
	fmt.Printf("%+v\n", params.ToplistItems)
	respondWithJSON(w, http.StatusOK, "ok")

}
