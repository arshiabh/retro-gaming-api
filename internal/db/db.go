package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/retry"
	_ "github.com/lib/pq"
)

func New(addr string, maxIdleConns, maxOpenConns int) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var db *sql.DB
	retry.Retry(ctx, func() error {
		var err error
		db, err = sql.Open("postgres", addr)
		if err != nil {
			return err
		}
		return db.Ping()

	}, retry.WithRetries(4))

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(time.Minute * 15)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Println("database connected succesfuly!")

	return db, nil
}

func WithTx(db *sql.DB, operation func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = operation(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
