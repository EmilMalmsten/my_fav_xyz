package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
)

func (cfg apiConfig) handlerToplistsGetOne(w http.ResponseWriter, r *http.Request) {
	toplistIDString := chi.URLParam(r, "toplistID")
	toplistID, err := strconv.Atoi(toplistIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	dbToplist, err := cfg.DB.GetToplist(toplistID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		respondWithError(w, http.StatusBadRequest, "Could not get toplist")
		return
	}

	respondWithJSON(w, http.StatusOK, dbToplist)
}
