package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
)

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Invalid user ID type")
		return
	}

	decoder := json.NewDecoder(r.Body)
	updateUserRequest := updateUserRequest{}
	err := decoder.Decode(&updateUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(updateUserRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	_, err = cfg.DB.UpdateUser(userID, updateUserRequest.Email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
