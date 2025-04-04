package main

import (
	"net/http"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/config"
	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	config *config.Config
	store  *store.Storage
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

		r.Get("/test" , app.HandleTest)
		r.Post("/users", app.HandleCreateUser)
		r.Get("/users", app.HandleCSRFTokenUser)

	})
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
