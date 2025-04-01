package main

import (
	"log"

	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

func main() {
	cfg := config.Load()
	store := store.NewStorage(nil)

	app := &application{
		config: *cfg,
		store:  store,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
