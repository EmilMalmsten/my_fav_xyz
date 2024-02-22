package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
		respondWithError(w, http.StatusInternalServerError, "Could not update toplist views")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

type ToplistLikeRequest struct {
	ToplistID int `json:"toplist_id"`
}

func (cfg *apiConfig) handlerToplistsLikes(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID type")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var toplistLikeRequest ToplistLikeRequest
	err := decoder.Decode(&toplistLikeRequest)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	err = cfg.DB.UpdateToplistLikes(toplistLikeRequest.ToplistID, userID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to like toplist")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
