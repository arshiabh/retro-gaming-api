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
	Set(*Score) (*Score, error)
	GetTopTen(int64) ([]*Score, error)
}

type PostgresScoreStore struct {
	db *sql.DB
}

func NewPostgresScoreStore(db *sql.DB) *PostgresScoreStore {
	return &PostgresScoreStore{
		db: db,
	}
}

func (s *PostgresScoreStore) Set(score *Score) (*Score, error) {
	query := `
	INSERT INTO scores (user_id, game_id, score) 
	VALUES ($1,$2,$3) RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, score.UserID, score.GameID, score.Point)
	if err := row.Scan(&score.ID); err != nil {
		return nil, err
	}
	return score, nil
}

func (s *PostgresScoreStore) GetTopTen(gameID int64) ([]*Score, error) {
	query := `
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	_, err := s.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
