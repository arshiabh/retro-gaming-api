package store

import (
	"database/sql"
	"time"
)

type score struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type ScoreStore interface{}

type PostgresScoreStore struct {
	db *sql.DB
}

func NewPostgresScoreStore(db *sql.DB) *PostgresScoreStore {
	return &PostgresScoreStore{
		db: db,
	}
}
