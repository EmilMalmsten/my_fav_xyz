package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbConfig struct {
	database *sql.DB
}

type Toplist struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ToplistItem struct {
	ID          int    `json:"id"`
	ListId      int    `json:"listId"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var ErrNotExist = errors.New("resource does not exist")

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

func (dbCfg *DbConfig) AddItemsToToplist(toplistItems []ToplistItem, listId int) error {

	query := "INSERT INTO list_items (toplist_id, rank, title, description) VALUES ($1, $2, $3, $4)"

	for i := range toplistItems {
		_, err := dbCfg.database.Exec(query,
			listId,
			toplistItems[i].Rank,
			toplistItems[i].Title,
			toplistItems[i].Description,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (dbCfg *DbConfig) UpdateToplist(toplist Toplist, listId int) error {
	var listsFound int
	doesToplistExist := "SELECT COUNT(*) FROM toplists WHERE id = $1"
	err := dbCfg.database.QueryRow(doesToplistExist, listId).Scan(&listsFound)
	if err != nil {
		return err
	}

	if listsFound < 1 {
		return ErrNotExist
	}

	query := "UPDATE toplists SET title = $1, description = $2 WHERE id = $3"
	_, err = dbCfg.database.Exec(query, toplist.Title, toplist.Description, listId)
	if err != nil {
		return err
	}

	return nil
}
