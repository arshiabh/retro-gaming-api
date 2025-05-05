package store

import "database/sql"

type Storage struct {
	Users  UserStore
	Games  GameStore
	Scores ScoreStore
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:  &PostgresUserStore{db: db},
		Games:  &PostgresGameStore{db: db},
		Scores: &PostgresScoreStore{db: db},
	}
}
