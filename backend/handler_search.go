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

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		http.Error(w, "search page is missing from request", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid search page param")
		return
	}

	resultLimit := 10
	offset := (page * resultLimit) - resultLimit
	searchResults, err := cfg.DB.SearchToplists(searchTerm, resultLimit, offset)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Toplist search failed")
		return
	}

	respondWithJSON(w, http.StatusOK, searchResults)
}