package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./task.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY,
		title TEXT,
		description TEXT,
		due_date TEXT,
		status TEXT
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
