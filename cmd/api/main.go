package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/auth"
	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/db"
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/ratelimiter"
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

	// for producer
	var pwg *sync.WaitGroup

	auth := auth.NewAuthentication(os.Getenv("secret_key"))
	kafka := kafka.NewKafkaService([]string{cfg.KafkaAddr})
	store := store.NewStorage(db)
	rdb := cache.NewStorage(cache.NewRedisClient("redis:6379"))

	service := service.NewService(store, kafka, rdb, pwg)

	errorLogger := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)

	ratelimiter := ratelimiter.NewFixedWindowLimiter(cfg.RateLimit.RequestPerTime, cfg.RateLimit.TimeFrame)

	app := &application{
		kafka:       kafka,
		config:      cfg,
		service:     service,
		auth:        auth,
		errorLogger: errorLogger,
		infoLogger:  infoLogger,
		ratelimiter: ratelimiter,
	}

	mux := app.mount()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// start the cleanup
	wg.Add(1)
	go ratelimiter.StartCleanup(ctx, &wg, time.Minute)

	// start service before shutdown to avoid race condition
	wg.Add(1)
	go kafka.StartConsumer(ctx, kafka.CreateReader("user-signup-consumer", "user-signup"), &wg)

	// listen for shutdown
	go func() {
		// order is important in this way the app got time to shutdown gracefully
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		//stop kafka
		pwg.Wait()
		kafka.Close()

		cancel()
		wg.Wait()
	}()

	if err := app.run(ctx, mux); err != nil {
		app.errorLogger.Fatal(err)
	}

}
