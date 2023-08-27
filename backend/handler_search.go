package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerToplistsSearch(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("term")
	if searchTerm == "" {
		http.Error(w, "Search term is missing", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		http.Error(w, "page size limit is missing from request", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid page limit param")
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		http.Error(w, "offset is missing from request", http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || limit < 1 {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid offset param")
		return
	}

	searchResults, err := cfg.DB.SearchToplists(searchTerm, limit, offset)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Toplist search failed")
		return
	}

	respondWithJSON(w, http.StatusOK, searchResults)
}