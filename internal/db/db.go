package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type ConfigDatabase struct {
	DataSource string `yaml:"data_source"`
}

func Connect(cfg *ConfigDatabase) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", cfg.DataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
