package store

import (
	"context"
	"database/sql"
	"time"
)

type UserStore interface {
	GetByUsername(string) (*User, error)
	Create(*User) (*User, error)
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) GetByUsername(username string) (*User, error) {
	query := `
	SELECT id, username, password_hash FROM users 
	WHERE username = ($1) 
	`
	user := &User{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresUserStore) Create(user *User) (*User, error) {
	query := `
	INSERT INTO users (username, password_hash, is_admin , created_at, updated_at) 
	VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, user.Username, user.Password, user.IsAdmin, user.CreatedAt, user.UpdatedAt)
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}
