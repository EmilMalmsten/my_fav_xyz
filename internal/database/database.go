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

func rankAlreadyExists(existingRanks []int, rank int) bool {
	for _, existingRank := range existingRanks {
		if existingRank == rank {
			return true
		}
	}
	return false
}

func (dbCfg *DbConfig) UpdateToplistItems(newListItems []ToplistItem, listId int) ([]ToplistItem, error) {

	existingListItems, err := dbCfg.GetToplistItems(listId)
	if err != nil {
		return []ToplistItem{}, err
	}

	if len(existingListItems) > len(newListItems) {
		lengthDifference := len(existingListItems) - len(newListItems)
		itemsToRemove := existingListItems[len(existingListItems)-lengthDifference:]
		err := dbCfg.DeleteSpecificToplistItems(itemsToRemove)
		if err != nil {
			return []ToplistItem{}, err
		}
	}

	var existingListItemRanks []int
	for _, item := range existingListItems {
		existingListItemRanks = append(existingListItemRanks, item.Rank)
	}

	var updatedItems []ToplistItem
	for i := range newListItems {
		var updatedItem ToplistItem
		if rankAlreadyExists(existingListItemRanks, newListItems[i].Rank) {
			updateQuery := "UPDATE list_items SET title = $1, description = $2 WHERE rank = $3 AND toplist_id = $4 RETURNING toplist_id, rank, title, description"
			err := dbCfg.database.QueryRow(updateQuery, newListItems[i].Title, newListItems[i].Description, newListItems[i].Rank, listId).Scan(
				&updatedItem.ListID,
				&updatedItem.Rank,
				&updatedItem.Title,
				&updatedItem.Description,
			)
			if err != nil {
				return []ToplistItem{}, err
			}

		} else {
			insertQuery := "INSERT INTO list_items (toplist_id, rank, title, description) VALUES ($1, $2, $3, $4) RETURNING toplist_id, rank, title, description"
			err := dbCfg.database.QueryRow(insertQuery, listId, newListItems[i].Rank, newListItems[i].Title, newListItems[i].Description).Scan(
				&updatedItem.ListID,
				&updatedItem.Rank,
				&updatedItem.Title,
				&updatedItem.Description,
			)
			if err != nil {
				return []ToplistItem{}, err
			}
		}
		updatedItems = append(updatedItems, updatedItem)
	}

	return updatedItems, nil
}

func (dbCfg *DbConfig) DeleteToplist(listId int) error {
	query := "DELETE FROM toplists WHERE id = $1"
	_, err := dbCfg.database.Exec(query, listId)
	if err != nil {
		return err
	}
	return nil
}

func (dbCfg *DbConfig) DeleteSpecificToplistItems(itemsToDelete []ToplistItem) error {
	for _, item := range itemsToDelete {
		query := "DELETE FROM list_items WHERE rank = $1"
		_, err := dbCfg.database.Exec(query, item.Rank)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dbCfg *DbConfig) DeleteToplistItems(listId int) error {
	query := "DELETE FROM list_items WHERE toplist_id = $1"
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
		if err == sql.ErrNoRows {
			return []ToplistItem{}, ErrNotExist
		}
		return []ToplistItem{}, err
	}
	defer rows.Close()

	var toplistItems []ToplistItem

	for rows.Next() {
		var item ToplistItem
		err := rows.Scan(&item.ListID, &item.Rank, &item.Title, &item.Description)
		if err != nil {
			fmt.Println(err)
			return []ToplistItem{}, err
		}
		item.ID = listId
		toplistItems = append(toplistItems, item)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return []ToplistItem{}, err
	}

	return toplistItems, nil
}
