package main

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update toplist")
		return

	}

	respondWithJSON(w, http.StatusOK, "")
}
