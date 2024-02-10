package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
)

type updateUserEmailRequest struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
	Password string `json:"password"`
}

type updateUserUsernameRequest struct {
	UserID      int    `json:"user_id"`
	NewUsername string `json:"new_username"`
	Password    string `json:"password"`
}

type updateUserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Email       string `json:"email"`
}

func (cfg *apiConfig) handlerUsersUpdateEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	updateUserEmailRequest := updateUserEmailRequest{}
	err := decoder.Decode(&updateUserEmailRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if !validEmail(updateUserEmailRequest.NewEmail) {
		respondWithError(w, http.StatusBadRequest, "Incorrect email format")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(updateUserEmailRequest.OldEmail)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(updateUserEmailRequest.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	updatedUser, err := cfg.DB.UpdateUserEmail(dbUser.ID, updateUserEmailRequest.NewEmail)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	updatedUserResp := database.User{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
	}

	respondWithJSON(w, http.StatusOK, updatedUserResp)
}

func (cfg *apiConfig) handlerUsersUpdateUsername(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	updateUserUsernameRequest := updateUserUsernameRequest{}
	err := decoder.Decode(&updateUserUsernameRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(updateUserUsernameRequest.NewUsername) < 3 {
		respondWithError(w, http.StatusBadRequest, "Username needs to be minimum 3 characters")
		return
	}

	dbUser, err := cfg.DB.GetUserByID(updateUserUsernameRequest.UserID)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(updateUserUsernameRequest.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	updatedUser, err := cfg.DB.UpdateUserUsername(dbUser.ID, updateUserUsernameRequest.NewUsername)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	updatedUserResp := database.User{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
		CreatedAt: updatedUser.CreatedAt,
	}

	respondWithJSON(w, http.StatusOK, updatedUserResp)
}

func (cfg *apiConfig) handlerUsersUpdatePassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	updateUserPasswordRequest := updateUserPasswordRequest{}
	err := decoder.Decode(&updateUserPasswordRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(updateUserPasswordRequest.NewPassword) < 8 {
		respondWithError(w, http.StatusBadRequest, "Password needs to be minimum 8 characters")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(updateUserPasswordRequest.Email)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(updateUserPasswordRequest.OldPassword, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	newPasswordHashed, err := auth.HashPassword(updateUserPasswordRequest.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	_, err = cfg.DB.UpdateUserPassword(dbUser.ID, newPasswordHashed)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
