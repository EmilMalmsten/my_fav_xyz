package main

import (
	"encoding/json"
	"net/http"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
)

type Toplist struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Items       []ToplistItem `json:"items"`
}

type ToplistItem struct {
	ID          int    `json:"id"`
	ListId      int    `json:"listId"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t Toplist) ToDBToplist() database.Toplist {
	dbItems := make([]database.ToplistItem, len(t.Items))
	for i, item := range t.Items {
		dbItems[i] = item.ToDBToplistItem()
	}

	return database.Toplist{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Items:       dbItems,
	}
}

func (t ToplistItem) ToDBToplistItem() database.ToplistItem {
	return database.ToplistItem{
		ID:          t.ID,
		ListId:      t.ListId,
		Rank:        t.Rank,
		Title:       t.Title,
		Description: t.Description,
	}
}

func (cfg apiConfig) handlerToplistsCreate(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Id int64 `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	var toplist Toplist
	err := decoder.Decode(&toplist)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	dbToplist := toplist.ToDBToplist()

	toplistId, err := cfg.DB.CreateToplist(dbToplist)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new toplist")
		return
	}

	respondWithJSON(w, http.StatusCreated, resp{
		Id: toplistId,
	})
}
