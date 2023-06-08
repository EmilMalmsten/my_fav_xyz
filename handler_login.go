package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/emilmalmsten/my_top_xyz/internal/auth"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginRequest := LoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(loginRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Hour,
		auth.TokenTypeAccess,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Hour*24*30,
		auth.TokenTypeRefresh,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
