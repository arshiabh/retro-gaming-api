package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct{}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	return r
}

func (app *application) run(max http.Handler) {
	srv := http.Server{
		Handler: max,
	}
	srv.ListenAndServe()
}
