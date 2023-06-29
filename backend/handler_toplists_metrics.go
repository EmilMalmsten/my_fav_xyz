package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
	"github.com/go-chi/chi"
)

func (cfg *apiConfig) handlerToplistsViews(w http.ResponseWriter, r *http.Request) {
	toplistIDString := chi.URLParam(r, "toplistID")

	toplistID, err := strconv.Atoi(toplistIDString)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	_, err = cfg.DB.UpdateToplistViews(toplistID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		respondWithError(w, http.StatusBadRequest, "Could not update toplist views")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
