package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arshiabh/retro-gaming-api/internal/auth"
	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/db"
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/service"
	"github.com/arshiabh/retro-gaming-api/internal/store/cache"

	"github.com/arshiabh/retro-gaming-api/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := db.New(cfg.DB.Addr, cfg.DB.MaxIdleConns, cfg.DB.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}

	auth := auth.NewAuthentication(os.Getenv("secret_key"))
	kafka := kafka.NewKafkaService([]string{cfg.KafkaAddr})
	store := store.NewStorage(db)
	rdb := cache.NewStorage(cache.NewRedisClient("localhost:6381"))

	service := service.NewService(store, kafka, rdb)

	app := &application{
		kafka:   kafka,
		config:  cfg,
		service: service,
		auth:    auth,
	}

	mux := app.mount()

	ctx, cancel := context.WithCancel(context.Background())
	go kafka.StartConsumer(ctx, kafka.CreateReader("user-signup-consumer", "user-signup"))

	go func() {
		if err := app.run(ctx, mux); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c
	//stop kafka
	cancel()
}
