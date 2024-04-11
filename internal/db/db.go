package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (db *sql.DB, err error) {
	dsn := os.Getenv("DB_URL")
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
