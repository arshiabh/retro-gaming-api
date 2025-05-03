package store

import (
	"context"
	"database/sql"
	"time"
)

type Game struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type GameStore interface {
	Create(*Game) (*Game, error)
}

type PostgresGameStore struct {
	db *sql.DB
}

func (s *PostgresGameStore) Create(game *Game) (*Game, error) {
	query := `
	INSERT INTO games (name, description)
	VALUES ($1,$2) RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, game.Name, game.Description)

	if err := row.Scan(&game.ID); err != nil {
		return nil, err
	}
	return game, nil
}
