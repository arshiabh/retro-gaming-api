package main

import (
	"log"
	"os"

	"github.com/arshiabh/retro-gaming-api/internal/auth"
	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/db"
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/service"

	"github.com/arshiabh/retro-gaming-api/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := db.New(cfg.DB.Addr, cfg.DB.MaxIdleConns, cfg.DB.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}

	auth := auth.NewAuthentication(os.Getenv("secret_key"))

	deps := module.Dependencies{
		Store: store.NewStorage(db),
		Kafka: kafka.NewClient([]string{"localhost:9092"}),
	}

	service := service.NewService(deps)

	app := &application{
		config:  cfg,
		service: service,
		auth:    auth,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
