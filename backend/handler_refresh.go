package main

import (
	"fmt"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
)

type refreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Refresh token missing in request auth header")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to check token revocation status")
		return
	}
	if isRevoked {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Refresh token is revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtSecret)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, refreshResponse{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Refresh token missing in request auth header")
		return
	}

	_, err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
