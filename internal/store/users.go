package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	GetByUserID()
	Create(*User) error
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) GetByUserID() {

}

func (s *PostgresUserStore) Create(user *User) error {
	query := `
	INSERT INTO users (username, password_hash, is_admin , created_at, updated_at) 
	VALUES ($1,$2,$3,$4,$5) 
	`
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
	return s.db.QueryRowContext(ctx, query, user.Username, string(hash), user.IsAdmin, user.CreatedAt, user.UpdatedAt).Err()
}
