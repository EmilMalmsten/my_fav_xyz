package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (cfg apiConfig) handlerToplistsDelete(w http.ResponseWriter, r *http.Request) {
	toplistIDString := chi.URLParam(r, "toplistID")
	toplistID, err := strconv.Atoi(toplistIDString)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	err = cfg.DB.DeleteToplist(toplistID)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete toplist")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
