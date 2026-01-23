package postgres

import "database/sql"

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{DB: db}
}
