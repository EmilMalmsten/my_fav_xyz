package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func (cfg *apiConfig) validateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
			return
		}
		userIDString, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
			return
		}

		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID")
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
