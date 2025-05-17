package store

import (
	"context"
	"database/sql"
	"time"
)

type Score struct {
	ID           int64     `json:"id,omitempty"`
	GameID       int64     `json:"game_id,omitempty"`
	UserID       int64     `json:"user_id,omitempty"`
	Point        int64     `json:"score,omitempty"`
	Submitted_at time.Time `json:"submitted_at,omitempty"`
	Created_at   time.Time `json:"-"`
	Updated_at   time.Time `json:"-"`
}

type LeaderBoard struct {
	Username     string    `db:"username" json:"username"`
	Score        int64     `db:"score" json:"score"`
	Submitted_at time.Time `db:"submitted_at" json:"submitted_at"`
}

type ScoreStore interface {
	Set(*Score) (*Score, error)
	GetTopTen(int64) ([]*LeaderBoard, error)
	GetUserScore(int64) ([]*Score, error)
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

func (s *PostgresScoreStore) GetTopTen(gameID int64) ([]*LeaderBoard, error) {
	query := `
	select u.username, s.score, s.submitted_at from scores s 
	inner join users u on u.id = s.user_id 
	where s.game_id = ($1)
	order by s.score
	DESC limit 10 
	`
	var result []*LeaderBoard

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &LeaderBoard{}
		if err := rows.Scan(&user.Username, &user.Score, &user.Submitted_at); err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func (s *PostgresScoreStore) GetUserScore(userID int64) ([]*Score, error) {
	query := `
	select game_id, score, submitted_at from scores 
	where user_id = ($1)
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	var scores []*Score

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		score := &Score{}
		rows.Scan(&score.GameID, &score.Point, &score.Submitted_at)

		scores = append(scores, score)
	}

	return scores, nil
}
