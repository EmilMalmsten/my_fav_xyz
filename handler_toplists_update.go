package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
)

func (cfg apiConfig) handlerToplistsUpdate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var toplist CreateToplistRequest
	err := decoder.Decode(&toplist)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	dbToplist := toplist.ToDBToplist()

	_, err = cfg.DB.UpdateToplist(dbToplist)
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
