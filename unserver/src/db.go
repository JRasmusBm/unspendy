package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func with_db(path string, run func(*sql.DB)) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	run(db)

	return nil
}
