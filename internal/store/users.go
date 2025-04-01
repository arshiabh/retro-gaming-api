package store

import "database/sql"

type UserStore interface {
	GetByUserID()
}

type User struct{}

type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) GetByUserID() {

}