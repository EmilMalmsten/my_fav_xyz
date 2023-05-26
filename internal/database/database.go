package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbConfig struct {
	database *sql.DB
}

type ToplistItem struct {
	ID          int    `json:"id"`
	ListId      int    `json:"listId"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func Init(dbUrl string) (*DbConfig, error) {
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return &DbConfig{}, err
	}
	err = db.Ping()
	if err != nil {
		return &DbConfig{}, err
	}

	return &DbConfig{database: db}, nil
}

func (dbCfg *DbConfig) CreateToplist(title string, description string) (int64, error) {
	var listId int64
	err := dbCfg.database.QueryRow(`
		INSERT INTO toplists (title, description)
		VALUES ($1, $2) RETURNING id`, title, description).Scan(&listId)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	return listId, nil
}

func (dbCfg *DbConfig) AddItemsToToplist(toplistItems []ToplistItem, listId int) ([]ToplistItem, error) {

	stmt, err := dbCfg.database.Prepare(`
		INSERT INTO list_items (toplist_id, rank, title, description)
		VALUES ($1, $2, $3, $4) RETURNING id
	`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()

	for i := range toplistItems {
		toplistItems[i].ListId = listId
		err := stmt.QueryRow(
			toplistItems[i].ListId,
			toplistItems[i].Rank,
			toplistItems[i].Title,
			toplistItems[i].Description,
		).Scan(&toplistItems[i].ID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return toplistItems, nil
}
