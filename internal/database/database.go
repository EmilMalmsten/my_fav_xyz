package database

import (
	"database/sql"

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

func (dbCfg *DbConfig) CreateToplist(title string) error {
	_, err := dbCfg.database.Exec("INSERT INTO lists (title) VALUES ($1)", title)
	if err != nil {
		return err
	}
	return nil
}
