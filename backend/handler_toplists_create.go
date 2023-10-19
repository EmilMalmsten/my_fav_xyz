package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
)

type toplistRequest struct {
	ToplistID   int                  `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	UserID      int                  `json:"user_id"`
	Items       []toplistItemRequest `json:"items"`
}

type toplistItemRequest struct {
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImagePath   string `json:"image_path"`
	Image       []byte `json:"-"`
}

func (t toplistRequest) ToDBToplist() database.Toplist {
	dbItems := make([]database.ToplistItem, len(t.Items))
	for i, item := range t.Items {
		dbItems[i] = item.ToDBToplistItem()
	}

	return database.Toplist{
		ToplistID:   t.ToplistID,
		Title:       t.Title,
		Description: t.Description,
		UserID:      t.UserID,
		Items:       dbItems,
	}
}

func (t toplistItemRequest) ToDBToplistItem() database.ToplistItem {
	return database.ToplistItem{
		Rank:        t.Rank,
		Title:       t.Title,
		Description: t.Description,
		ImagePath:   t.ImagePath,
		Image:       t.Image,
	}
}

func (cfg *apiConfig) handlerToplistsCreate(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Id int `json:"id"`
	}

	userIDValue := r.Context().Value(userIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID data type")
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

	err = validateToplistValues(toplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbToplist := toplist.ToDBToplist()
	insertedToplist, err := cfg.DB.InsertToplist(dbToplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new toplist")
		return
	}

	respondWithJSON(w, http.StatusCreated, resp{
		Id: insertedToplist.ToplistID,
	})
}

func validateToplistValues(toplist toplistRequest) error {
	if toplist.Title == "" {
		return errors.New("toplist title is missing")
	}
	maxDescriptionLength := 1000
	if len(toplist.Description) > maxDescriptionLength {
		return fmt.Errorf(
			"toplist description too long. Max is %d characters, got %d characters",
			maxDescriptionLength,
			len(toplist.Description),
		)
	}

	err := validateListItemValues(toplist.Items)
	if err != nil {
		return err
	}

	return nil
}

func validateListItemValues(toplistItems []toplistItemRequest) error {
	itemRanks := make([]int, len(toplistItems))
	for i, item := range toplistItems {
		itemRanks[i] = item.Rank
		if item.Title == "" {
			return errors.New("title missing in toplist items")
		}
	}

	sort.Ints(itemRanks)

	for i := 0; i < len(itemRanks); i++ {
		if itemRanks[i] != i+1 {
			return errors.New("toplist item ranks are not ordered properly")
		}
	}
	return nil
}
