package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
)

func (cfg *apiConfig) handlerToplistsUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse form data")
		return
	}



	toplistIDStr := r.FormValue("id")
	toplistID, err := strconv.Atoi(toplistIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid toplist ID")
		return
	}
	title := r.FormValue("title")
	description := r.FormValue("description")
	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Invalid user ID type")
		return
	}

	items := make([]toplistItemRequest, 0)
	for index := 0; ; index++ {
		key := fmt.Sprintf("items[%d][", index)

		// Check if the title key exists in the form data
		_, hasTitle := r.Form[key+"title]"]
		if !hasTitle {
			// Key not found, exit the loop
			break
		}

		item := toplistItemRequest{
			Title:       r.FormValue(key + "title]"),
			Description: r.FormValue(key + "description]"),
			ImagePath: 	 r.FormValue(key + "path]"),
			Rank:        -1, // Set a default value for Rank
		}

		rankStr := r.FormValue(key + "rank]")
		if rankStr != "" {
			rank, err := strconv.Atoi(rankStr)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid item rank")
				return
			}
			item.Rank = rank
		}

		// Check if the image key exists in the form data
		_, fileHeader, err := r.FormFile(key + "image]")
		if err == nil && fileHeader != nil {
			file, err := fileHeader.Open()
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to read item image")
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to read item image")
				return
			}

			item.Image = fileBytes
		}

		items = append(items, item)
	}

	toplist := toplistRequest{
		ToplistID:   toplistID,
		Title:       title,
		Description: description,
		UserID:      userID,
		Items:       items,
	}

	err = validateToplistValues(toplist)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbToplist := toplist.ToDBToplist()

	_, err = cfg.DB.UpdateToplist(dbToplist)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update toplist")
		return

	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
