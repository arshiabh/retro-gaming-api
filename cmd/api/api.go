package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct {
	cfg config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	return r
}

func (app *application) run(max http.Handler) {
	srv := http.Server{
		Addr:    app.cfg.addr,
		Handler: max,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
