package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type List struct {
	Title string
}

func Db(dbUrl string) {
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var title string
	err = db.QueryRow("SELECT title FROM lists LIMIT 1").Scan(&title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("found title: %+v\n", title)
}
