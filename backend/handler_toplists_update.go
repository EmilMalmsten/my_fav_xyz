package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
)

func (cfg *apiConfig) handlerToplistsUpdate(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID type")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var toplist toplistRequest
	err := decoder.Decode(&toplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	toplist.UserID = userID

	dbToplist := toplist.ToDBToplist()

	_, err = cfg.DB.UpdateToplist(dbToplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update toplist")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg *apiConfig) handlerToplistsUpdateItems(w http.ResponseWriter, r *http.Request) {
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

		_, hasTitle := r.Form[key+"title]"]
		if !hasTitle {
			break
		}

		item := toplistItemRequest{
			Title:       r.FormValue(key + "title]"),
			Description: r.FormValue(key + "description]"),
			ImagePath:   r.FormValue(key + "path]"),
			Rank:        -1,
		}

		rankStr := r.FormValue(key + "rank]")
		if rankStr == "" {
			respondWithError(w, http.StatusInternalServerError, "Missing item rank")
			return
		}
		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Item rank string conversion failed")
			return
		}
		item.Rank = rank

		imageBase64 := r.FormValue(key + "image]")
		if imageBase64 != "" {
			dataParts := strings.SplitN(imageBase64, ",", 2)
			if len(dataParts) != 2 {
				respondWithError(w, http.StatusBadRequest, "Invalid base64 image data")
				return
			}
			base64Data := dataParts[1]

			imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to decode base64 image")
				return
			}

			item.Image = imageBytes
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

	_, err = cfg.DB.UpdateToplistItems(dbToplist.Items, dbToplist.ToplistID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Toplist does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update toplist")
		return

	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
