package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
)

func (cfg apiConfig) handlerToplistsUpdate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var toplist database.Toplist
	err := decoder.Decode(&toplist)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	listIdString := chi.URLParam(r, "listId")
	listId, err := strconv.Atoi(listIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	err = cfg.DB.UpdateToplist(toplist, listId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new toplist")
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
