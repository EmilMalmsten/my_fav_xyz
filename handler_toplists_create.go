package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
)

type createToplistRequest struct {
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	UserID      int                        `json:"user_id"`
	Items       []createToplistItemRequest `json:"items"`
}

type createToplistItemRequest struct {
	ListId      int    `json:"listId"`
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
		ListID:      t.ListId,
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
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
		Id: insertedToplist.ID,
	})
}
