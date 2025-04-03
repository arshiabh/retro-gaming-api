package main

import (
	"log"
	"net/http"

	"github.com/arshiabh/retro-gaming-api/internal/store"
)

func (app *application) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	var users []store.User
	writeJSON(w, http.StatusOK, users)
}

func (app *application) HandleCSRFTokenUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, "message: token set")
}

func (app *application) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := &store.User{}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}
	user.Username = username
	user.Password = password
	if err := app.store.Users.Create(user); err != nil {
		log.Println(err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, "message: user successfully created!")
}
