package main

import (
	"log"
	"os"

	"github.com/arshiabh/retro-gaming-api/internal/auth"
	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/db"
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/service"

	// "github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := db.New(cfg.DB.Addr, cfg.DB.MaxIdleConns, cfg.DB.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}

	// store := store.NewStorage(db)
	auth := auth.NewAuthentication(os.Getenv("secret_key"))
	// kafka := kafka.NewClient([]string{"localhost:9999"})

	deps := module.Dependencies{
		Store: store.NewStorage(db),
	}

	service := service.NewService(deps)

	app := &application{
		config:  cfg,
		service: service,
		// store:  store,
		auth: auth,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
