package store

import "database/sql"

type UserStore interface {
	GetByUserID()
}

type User struct{}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) GetByUserID() {}
