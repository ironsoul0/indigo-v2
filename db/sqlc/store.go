package db

import "database/sql"

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
