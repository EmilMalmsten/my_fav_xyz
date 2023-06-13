package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
)

func (cfg *apiConfig) handlerToplistsGetOne(w http.ResponseWriter, r *http.Request) {
	toplistIDString := chi.URLParam(r, "toplistID")
	toplistID, err := strconv.Atoi(toplistIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}

	dbToplist, err := cfg.DB.GetToplist(toplistID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		respondWithError(w, http.StatusBadRequest, "Could not get toplist")
		return
	}

	respondWithJSON(w, http.StatusOK, dbToplist)
}

func (cfg apiConfig) handlerToplistsGetMany(w http.ResponseWriter, r *http.Request) {
	pageIDString := r.URL.Query().Get("page_id")
	pageID, err := strconv.Atoi(pageIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid page ID parameter")
		return
	}
	if pageID < 1 {
		respondWithError(w, http.StatusBadRequest, "Page ID value needs to be minimum 1")
		return
	}

	pageSizeString := r.URL.Query().Get("page_size")
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid page size parameter")
		return
	}
	maxPageSize := 20
	if pageSize < 1 || pageSize > maxPageSize {
		respondWithError(w, http.StatusBadRequest, "Page size value needs to be minimum 1")
		return
	}

	offset := (pageID - 1) * pageSize
	toplists, err := cfg.DB.ListToplists(pageSize, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get toplists")
		return
	}

	respondWithJSON(w, http.StatusOK, toplists)
}
