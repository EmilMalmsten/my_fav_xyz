package database

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrNotExist = errors.New("resource does not exist")
var ErrAlreadyExist = errors.New("already exists")
var ErrIsExpired = errors.New("already expired")

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

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
