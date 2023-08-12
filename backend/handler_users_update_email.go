package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
)

type updateUserEmailRequest struct {
	OldEmail    string `json:"old_email"`
	NewEmail    string `json:"new_email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersUpdateEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	updateUserEmailRequest := updateUserEmailRequest{}
	err := decoder.Decode(&updateUserEmailRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
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
		ID: updatedUser.ID,
		Email: updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
	}

	fmt.Println(updatedUserResp)

	respondWithJSON(w, http.StatusOK, updatedUserResp)
}
