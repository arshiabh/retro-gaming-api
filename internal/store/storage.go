package store

import "database/sql"

type Storage struct {
	Users UserStore
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users: &PostgresUserStore{db: db},
	}
}
