package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func (dbCfg *DbConfig) InsertToplist(toplist Toplist) (Toplist, error) {

	insertQuery := `
		INSERT INTO toplists (title, description, user_id)
		VALUES ($1, $2, $3) 
		RETURNING id, title, description, user_id, created_at
	`

	var insertedToplist Toplist

	err := dbCfg.database.QueryRow(insertQuery, toplist.Title, toplist.Description, toplist.UserID).Scan(
		&insertedToplist.ToplistID,
		&insertedToplist.Title,
		&insertedToplist.Description,
		&insertedToplist.UserID,
		&insertedToplist.CreatedAt,
	)
	if err != nil {
		return insertedToplist, err
	}

	insertedListItems, err := dbCfg.InsertToplistItems(toplist.Items, insertedToplist.ToplistID)
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
		RETURNING id, toplist_id, rank, title, description`)
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
			&insertedItem.ItemID,
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
	row := dbCfg.database.QueryRowContext(context.Background(), query, toplist.ToplistID)

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
		RETURNING id, title, description, user_id, created_at
	`

	var updatedToplist Toplist

	err = dbCfg.database.QueryRow(updateQuery, toplist.Title, toplist.Description, toplist.ToplistID).Scan(
		&updatedToplist.ToplistID,
		&updatedToplist.Title,
		&updatedToplist.Description,
		&updatedToplist.UserID,
		&updatedToplist.CreatedAt,
	)
	if err != nil {
		return Toplist{}, err
	}

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

	for i := range newListItems {
		dbCfg.saveOrDeleteImage(&newListItems[i], listId)
		log.Println(newListItems[i].ImagePath)
	}

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
			updateQuery := `
				UPDATE list_items SET title = $1, description = $2, image_path = $3
				WHERE rank = $4 AND toplist_id = $5
				RETURNING id, toplist_id, rank, title, description, image_path
			`
			err := dbCfg.database.QueryRow(
				updateQuery,
				newListItems[i].Title,
				newListItems[i].Description,
				newListItems[i].ImagePath,
				newListItems[i].Rank,
				listId,
			).Scan(
				&updatedItem.ItemID,
				&updatedItem.ListID,
				&updatedItem.Rank,
				&updatedItem.Title,
				&updatedItem.Description,
				&updatedItem.ImagePath,
			)
			if err != nil {
				return []ToplistItem{}, err
			}

		} else {
			insertQuery := `
				INSERT INTO list_items (toplist_id, rank, title, description, image_path)
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING id, toplist_id, rank, title, description, image_path
			`
			err := dbCfg.database.QueryRow(
				insertQuery,
				listId,
				newListItems[i].Rank,
				newListItems[i].Title,
				newListItems[i].Description,
				newListItems[i].ImagePath,
			).Scan(
				&updatedItem.ItemID,
				&updatedItem.ListID,
				&updatedItem.Rank,
				&updatedItem.Title,
				&updatedItem.Description,
				&updatedItem.ImagePath,
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

	err = dbCfg.DeleteToplistItems(listId)
	if err != nil {
		return err
	}

	err = deleteToplistImages(listId)
	if err != nil {
		return err
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

func (dbCfg *DbConfig) GetToplistLikes(toplistID int) ([]int, error) {
	query := "SELECT user_id FROM toplist_likes WHERE toplist_id = $1"
	rows, err := dbCfg.database.Query(query, toplistID)
	if err != nil {
		return nil, err
	}

	var userIDs []int
	for rows.Next() {
		var userID int
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func (dbCfg *DbConfig) GetToplist(listId int) (Toplist, error) {
	var toplist Toplist

	query := "SELECT id, title, description, user_id, created_at FROM toplists WHERE id = $1"
	row := dbCfg.database.QueryRowContext(context.Background(), query, listId)

	err := row.Scan(
		&toplist.ToplistID,
		&toplist.Title,
		&toplist.Description,
		&toplist.UserID,
		&toplist.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return toplist, ErrNotExist
		}
		return toplist, err
	}

	user, err := dbCfg.GetUserByID(toplist.UserID)
	if err != nil {
		return toplist, err
	}
	toplist.Username = user.Username

	toplistItems, err := dbCfg.GetToplistItems(listId)
	if err != nil {
		fmt.Println(err)
		return toplist, err
	}

	toplistLikes, err := dbCfg.GetToplistLikes(listId)
	if err != nil {
		return toplist, err
	}
	toplist.LikeIDs = toplistLikes
	toplist.LikeCount = len(toplistLikes)

	toplist.Items = toplistItems
	return toplist, nil
}

func (dbCfg *DbConfig) GetToplistItems(listId int) ([]ToplistItem, error) {
	query := "SELECT id, toplist_id, rank, title, description, image_path FROM list_items WHERE toplist_id = $1"
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
		err := rows.Scan(&item.ItemID, &item.ListID, &item.Rank, &item.Title, &item.Description, &item.ImagePath)
		if err != nil {
			fmt.Println(err)
			return []ToplistItem{}, err
		}
		toplistItems = append(toplistItems, item)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return []ToplistItem{}, err
	}

	return toplistItems, nil
}

func (dbCfg *DbConfig) ListToplists(limit, offset int) ([]Toplist, error) {
	query := `
		SELECT id, title, description, user_id, created_at FROM toplists 
		ORDER BY id
		LIMIT $1 
		OFFSET $2
	`
	rows, err := dbCfg.database.QueryContext(context.Background(), query, limit, offset)
	if err != nil {
		return []Toplist{}, err
	}
	defer rows.Close()

	var toplists []Toplist
	for rows.Next() {
		var toplist Toplist
		err := rows.Scan(
			&toplist.ToplistID,
			&toplist.Title,
			&toplist.Description,
			&toplist.UserID,
			&toplist.CreatedAt,
		)
		if err != nil {
			return []Toplist{}, err
		}
		toplist.Items, err = dbCfg.GetToplistItems(toplist.ToplistID)
		if err != nil {
			return []Toplist{}, err
		}
		toplists = append(toplists, toplist)
	}

	return toplists, nil
}

func (dbCfg *DbConfig) ListToplistsByProperty(limit int, property string) ([]Toplist, error) {
	orderClause := ""

	switch property {
	case "date":
		orderClause = "created_at DESC"
	case "views":
		orderClause = "views DESC"
	case "likes":
		orderClause = "likes DESC"
	default:
		orderClause = "created_at DESC"
	}

	query := ` 
		SELECT id, title, description, user_id, created_at, views FROM toplists 
		ORDER BY ` + orderClause + `
		LIMIT $1
	`

	rows, err := dbCfg.database.QueryContext(context.Background(), query, limit)
	if err != nil {
		return []Toplist{}, err
	}
	defer rows.Close()

	var toplists []Toplist
	for rows.Next() {
		var toplist Toplist
		err := rows.Scan(
			&toplist.ToplistID,
			&toplist.Title,
			&toplist.Description,
			&toplist.UserID,
			&toplist.CreatedAt,
			&toplist.Views,
		)
		if err != nil {
			return []Toplist{}, err
		}
		user, err := dbCfg.GetUserByID(toplist.UserID)
		if err != nil {
			return []Toplist{}, err
		}
		toplist.Username = user.Username

		toplistItems, err := dbCfg.GetToplistItems(toplist.ToplistID)
		if err != nil {
			return []Toplist{}, err
		}

		toplist.Items = toplistItems
		toplists = append(toplists, toplist)
	}

	return toplists, nil
}

func (dbCfg *DbConfig) UpdateToplistViews(toplistID int) (Toplist, error) {

	updateQuery := `
		UPDATE toplists SET views = views + 1
		WHERE id = $1
		RETURNING id, title, description, user_id, created_at, views
	`

	var updatedToplist Toplist

	err := dbCfg.database.QueryRow(updateQuery, toplistID).Scan(
		&updatedToplist.ToplistID,
		&updatedToplist.Title,
		&updatedToplist.Description,
		&updatedToplist.UserID,
		&updatedToplist.CreatedAt,
		&updatedToplist.Views,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Toplist{}, ErrNotExist
		} else {
			return Toplist{}, err
		}
	}
	return updatedToplist, nil
}

func (dbCfg *DbConfig) UpdateToplistLikes(toplistID, userID int) error {

	var exists bool
	err := dbCfg.database.QueryRow("SELECT EXISTS(SELECT 1 FROM toplist_likes WHERE user_id = $1 AND toplist_id = $2)",
		userID,
		toplistID,
	).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if !exists {
		insertQuery := "INSERT INTO toplist_likes (user_id, toplist_id, liked_at) VALUES ($1, $2, CURRENT_TIMESTAMP)"
		_, err = dbCfg.database.Exec(insertQuery, userID, toplistID)
		if err != nil {
			return err
		}
		return nil
	}

	deleteQuery := "DELETE FROM toplist_likes WHERE user_id = $1 AND toplist_id = $2"
	_, err = dbCfg.database.Exec(deleteQuery, userID, toplistID)
	if err != nil {
		return err
	}
	return nil
}

func (dbCfg *DbConfig) SearchToplists(searchTerm string, limit, offset int) ([]Toplist, error) {
	query := `
		SELECT id, title, description, user_id, created_at FROM toplists
		WHERE title ILIKE '%' || $1 || '%'
		ORDER BY id desc
		LIMIT $2
		OFFSET $3
	`

	rows, err := dbCfg.database.QueryContext(context.Background(), query, searchTerm, limit, offset)
	if err != nil {
		return []Toplist{}, err
	}
	defer rows.Close()

	var toplists []Toplist
	for rows.Next() {
		var toplist Toplist
		err := rows.Scan(
			&toplist.ToplistID,
			&toplist.Title,
			&toplist.Description,
			&toplist.UserID,
			&toplist.CreatedAt,
		)
		if err != nil {
			return []Toplist{}, err
		}
		toplists = append(toplists, toplist)
	}

	return toplists, nil
}

func (dbCfg *DbConfig) ListToplistsByUser(userID, limit, offset int) ([]Toplist, error) {
	query := `
		SELECT id, title, description, user_id, created_at FROM toplists
		WHERE user_id = $1
		ORDER BY created_at desc
		LIMIT $2
		OFFSET $3
	`

	rows, err := dbCfg.database.QueryContext(context.Background(), query, userID, limit, offset)
	if err != nil {
		return []Toplist{}, err
	}
	defer rows.Close()

	var toplists []Toplist
	for rows.Next() {
		var toplist Toplist
		err := rows.Scan(
			&toplist.ToplistID,
			&toplist.Title,
			&toplist.Description,
			&toplist.UserID,
			&toplist.CreatedAt,
		)
		if err != nil {
			return []Toplist{}, err
		}
		user, err := dbCfg.GetUserByID(toplist.UserID)
		if err != nil {
			return []Toplist{}, err
		}
		toplist.Username = user.Username

		toplistItems, err := dbCfg.GetToplistItems(toplist.ToplistID)
		if err != nil {
			return []Toplist{}, err
		}

		toplist.Items = toplistItems

		toplists = append(toplists, toplist)
	}

	return toplists, nil
}
