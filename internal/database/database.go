package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbConfig struct {
	database *sql.DB
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
	err := dbCfg.database.QueryRow("INSERT INTO toplists (title, description) VALUES ($1, $2) RETURNING id", title, description).Scan(&listId)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	return listId, nil
}
