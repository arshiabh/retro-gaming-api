package main

import (
	"log"

	"github.com/arshiabh/retro-gaming-api/internal/db"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

func main() {
	DB, err := db.New("host=localhost port=5555 user=admin password=secret dbname=retro-game sslmode=disable", 30, 30)
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	store.NewStorage(DB)
}
