package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
)

func (cfg *apiConfig) handlerToplistsDelete(w http.ResponseWriter, r *http.Request) {
	toplistIDString := chi.URLParam(r, "toplistID")

	toplistID, err := strconv.Atoi(toplistIDString)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID type")
		return
	}

	dbToplist, err := cfg.DB.GetToplist(toplistID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not find toplist")
		return
	}

	if dbToplist.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized to remove toplist")
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
