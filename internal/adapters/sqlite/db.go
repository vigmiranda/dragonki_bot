package sqlite

import "database/sql"

type SQLite struct {
	DB *sql.DB
}

func NewSQLite(db *sql.DB) *SQLite {
	return &SQLite{DB: db}
}
