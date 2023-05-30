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
var ErrAlreadyExist = errors.New("already exists")

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

func (dbCfg *DbConfig) getExistingListItemRanks(listId int) ([]int, error) {
	query := "SELECT rank FROM list_items WHERE toplist_id = $1"
	rows, err := dbCfg.database.Query(query, listId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var existingRanks []int
	for rows.Next() {
		var rank int
		err := rows.Scan(&rank)
		if err != nil {
			return nil, err
		}
		existingRanks = append(existingRanks, rank)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return existingRanks, nil
}

func rankExists(existingRanks []int, rank int) bool {
	for _, existingRank := range existingRanks {
		if existingRank == rank {
			return true
		}
	}
	return false
}

func (dbCfg *DbConfig) AddItemsToToplist(toplistItems []ToplistItem, listId int) error {
	existingListItemRanks, err := dbCfg.getExistingListItemRanks(listId)
	if err != nil {
		return err
	}
	fmt.Println(existingListItemRanks)

	for i := range toplistItems {
		if rankExists(existingListItemRanks, toplistItems[i].Rank) {
			updateQuery := "UPDATE list_items SET title = $1, description = $2 WHERE rank = $3 AND toplist_id = $4"
			_, err := dbCfg.database.Exec(updateQuery, toplistItems[i].Title, toplistItems[i].Description, toplistItems[i].Rank, listId)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			insertQuery := "INSERT INTO list_items (toplist_id, rank, title, description) VALUES ($1, $2, $3, $4)"
			_, err := dbCfg.database.Exec(insertQuery, listId, toplistItems[i].Rank, toplistItems[i].Title, toplistItems[i].Description)
			if err != nil {
				fmt.Println(err)
				return err
			}
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
