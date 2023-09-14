package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (cfg *apiConfig) handlerUsersGetByID(w http.ResponseWriter, r *http.Request) {
	type getUserResp struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}

	userIDString := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	dbUser, err := cfg.DB.GetUserByID(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	respondWithJSON(w, http.StatusOK, getUserResp{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
	})
}
