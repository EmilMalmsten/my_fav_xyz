package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
)

type createToplistRequest struct {
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	UserID      int                        `json:"user_id"`
	Items       []createToplistItemRequest `json:"items"`
}

type createToplistItemRequest struct {
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t createToplistRequest) ToDBToplist() database.Toplist {
	dbItems := make([]database.ToplistItem, len(t.Items))
	for i, item := range t.Items {
		dbItems[i] = item.ToDBToplistItem()
	}

	return database.Toplist{
		Title:       t.Title,
		Description: t.Description,
		UserID:      t.UserID,
		Items:       dbItems,
	}
}

func (t createToplistItemRequest) ToDBToplistItem() database.ToplistItem {
	return database.ToplistItem{
		Rank:        t.Rank,
		Title:       t.Title,
		Description: t.Description,
	}
}

func (cfg apiConfig) handlerToplistsCreate(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Id int `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	var toplist createToplistRequest
	err := decoder.Decode(&toplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	err = validateToplist(toplist)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	dbToplist := toplist.ToDBToplist()

	insertedToplist, err := cfg.DB.InsertToplist(dbToplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new toplist")
		return
	}

	respondWithJSON(w, http.StatusCreated, resp{
		Id: insertedToplist.ID,
	})
}

func validateToplist(toplist createToplistRequest) error {
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
	// TODO: validate whole toplist item, not just rank

	ok := validateItemRanks(toplist.Items)
	if !ok {
		return errors.New("toplist item ranks are not ordered properly")
	}

	return nil
}

func validateItemRanks(toplistItems []createToplistItemRequest) bool {
	itemRanks := make([]int, len(toplistItems))
	for i, item := range toplistItems {
		itemRanks[i] = item.Rank
	}

	sort.Ints(itemRanks)

	for i := 0; i < len(itemRanks); i++ {
		if itemRanks[i] != i+1 {
			return false
		}
	}
	return true
}
