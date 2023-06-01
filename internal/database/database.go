package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbConfig struct {
	database *sql.DB
}

type Toplist struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Items       []ToplistItem `json:"items"`
}

type ToplistItem struct {
	ID          int    `json:"id"`
	ListID      int    `json:"listId"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var ErrNotExist = errors.New("resource does not exist")
var ErrAlreadyExist = errors.New("already exists")

func CreateDatabaseConnection(dbUrl string) (*DbConfig, error) {
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

func (dbCfg *DbConfig) InsertToplist(toplist Toplist) (Toplist, error) {

	insertQuery := `
		INSERT INTO toplists (title, description)
		VALUES ($1, $2) 
		RETURNING id, title, description
	`

	var insertedToplist Toplist

	err := dbCfg.database.QueryRow(insertQuery, toplist.Title, toplist.Description).Scan(
		&insertedToplist.ID,
		&insertedToplist.Title,
		&insertedToplist.Description,
	)
	if err != nil {
		return insertedToplist, err
	}

	insertedListItems, err := dbCfg.InsertToplistItems(toplist.Items, insertedToplist.ID)
	if err != nil {
		fmt.Println(err)
		return insertedToplist, err
	}

	insertedToplist.Items = insertedListItems

	return insertedToplist, nil
}

func (dbCfg *DbConfig) InsertToplistItems(toplistItems []ToplistItem, listId int) ([]ToplistItem, error) {
	tx, err := dbCfg.database.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO list_items(toplist_id, rank, title, description) 
		VALUES ($1, $2, $3, $4) 
		RETURNING toplist_id, rank, title, description`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var insertedListItems []ToplistItem

	for _, item := range toplistItems {
		var insertedItem ToplistItem

		err := stmt.QueryRow(
			listId,
			item.Rank,
			item.Title,
			item.Description,
		).Scan(
			&insertedItem.ListID,
			&insertedItem.Rank,
			&insertedItem.Title,
			&insertedItem.Description,
		)
		if err != nil {
			return nil, err
		}

		insertedListItems = append(insertedListItems, insertedItem)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return insertedListItems, nil
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

func (dbCfg *DbConfig) UpdateToplist(toplist Toplist) (Toplist, error) {

	query := "SELECT 1 FROM toplists WHERE id = $1"
	row := dbCfg.database.QueryRowContext(context.Background(), query, toplist.ID)

	var toplistExists bool
	err := row.Scan(&toplistExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return Toplist{}, ErrNotExist
		} else {
			return Toplist{}, err
		}
	}

	updateQuery := `
		UPDATE toplists SET title = $1, description = $2
		WHERE id = $3
		RETURNING id, title, description
	`

	var updatedToplist Toplist

	err = dbCfg.database.QueryRow(updateQuery, toplist.Title, toplist.Description, toplist.ID).Scan(
		&updatedToplist.ID,
		&updatedToplist.Title,
		&updatedToplist.Description,
	)
	if err != nil {
		return Toplist{}, err
	}

	updatedItems, err := dbCfg.UpdateToplistItems(toplist.Items, toplist.ID)
	if err != nil {
		return Toplist{}, err
	}

	updatedToplist.Items = updatedItems
	return updatedToplist, nil
}

func (dbCfg *DbConfig) UpdateToplistItems(toplistItems []ToplistItem, listId int) ([]ToplistItem, error) {

	err := dbCfg.DeleteToplistItems(listId)
	if err != nil {
		return nil, err
	}

	newItems, err := dbCfg.InsertToplistItems(toplistItems, listId)
	if err != nil {
		return nil, err
	}

	return newItems, nil
}

func (dbCfg *DbConfig) DeleteToplistItems(listId int) error {
	query := "DELETE FROM list_items WHERE toplist_id = $1"
	_, err := dbCfg.database.Exec(query, listId)
	if err != nil {
		return err
	}
	return nil
}

func (dbCfg *DbConfig) DeleteToplist(listId int) error {
	query := "DELETE FROM toplists WHERE id = $1"
	_, err := dbCfg.database.Exec(query, listId)
	if err != nil {
		return err
	}
	return nil
}

func (dbCfg *DbConfig) GetToplist(listId int) (Toplist, error) {
	var toplist Toplist

	query := "SELECT id, title, description FROM toplists WHERE id = $1"
	row := dbCfg.database.QueryRowContext(context.Background(), query, listId)

	err := row.Scan(&toplist.ID, &toplist.Title, &toplist.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return toplist, ErrNotExist
		}
		return toplist, err
	}

	toplistItems, err := dbCfg.GetToplistItems(listId)
	if err != nil {
		fmt.Println(err)
		return toplist, err
	}

	toplist.Items = toplistItems
	return toplist, nil
}

func (dbCfg *DbConfig) GetToplistItems(listId int) ([]ToplistItem, error) {
	query := "SELECT toplist_id, rank, title, description FROM list_items WHERE toplist_id = $1"
	rows, err := dbCfg.database.QueryContext(context.Background(), query, listId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var toplistItems []ToplistItem

	for rows.Next() {
		var item ToplistItem
		err := rows.Scan(&item.ListID, &item.Rank, &item.Title, &item.Description)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		item.ID = listId
		toplistItems = append(toplistItems, item)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return toplistItems, nil
}
