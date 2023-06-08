package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/auth"
)

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	updateUserRequest := updateUserRequest{}
	err = decoder.Decode(&updateUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(updateUserRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	_, err = cfg.DB.UpdateUser(userIDInt, updateUserRequest.Email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
