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

func WithTx(db *sql.DB, operation func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := operation(tx); err != nil {
		_ = tx.Rollback()
	}

	return tx.Commit()
}
