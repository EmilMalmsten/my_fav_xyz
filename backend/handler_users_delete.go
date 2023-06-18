package main

import "net/http"

func (cfg *apiConfig) handlerUsersDelete(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Invalid user ID type")
		return
	}

	err := cfg.DB.DeleteUser(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete user")
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
