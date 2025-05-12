package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/auth"
	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	kafka   *kafka.KafkaService
	config  *config.Config
	service *service.Service
	auth    auth.Authenticator
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(csrfMiddleware)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		//check for only specific origin to avoid misconfig cors
		AllowedOrigins:   []string{"https://*"},
		AllowCredentials: true,
	}))

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/login", app.HandleLoginUser)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.HandleCreateUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.JWTAuthMiddleware)
			r.Post("/games", app.HandleCreateGame)
			r.Post("/scores", app.HandleSetScore)
		})

		r.Get("/test", app.HandleTest)

	})
	return r
}

func (app *application) run(ctx context.Context, mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	serverErr := make(chan error, 1)

	go func() {
		//go in background run server if we got error catch it with chan
		log.Printf("starting http server on port %s\n", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	//waiting for either of them to close, error or done
	select {
	case <-ctx.Done():
		log.Println("shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println("error while shutting down server")
			return err
		}
		return nil

	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
