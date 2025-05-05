package store

import (
	"context"
	"database/sql"
	"time"
)

type Score struct {
	ID           int64     `json:"id"`
	GameID       int64     `json:"game_id"`
	UserID       int64     `json:"user_id"`
	Point        int64     `json:"score"`
	Submitted_at time.Time `json:"submitted_at"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

type ScoreStore interface {
	Set(int64, int64) (*Score, error)
}

type PostgresScoreStore struct {
	db *sql.DB
}

func NewPostgresScoreStore(db *sql.DB) *PostgresScoreStore {
	return &PostgresScoreStore{
		db: db,
	}
}

func (s *PostgresScoreStore) Set(gameID, point int64) (*Score, error) {
	query := `
	INSERT INTO scores VALUES ()
	`
	score := &Score{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, gameID, point)
	if err := row.Scan(); err != nil {
		return nil, err
	}
	return score, nil
}
