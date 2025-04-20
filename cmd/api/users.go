package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arshiabh/retro-gaming-api/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type userPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	var users []store.User
	writeJSON(w, http.StatusOK, users)
}

func (app *application) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload userPayload
	if err := readJSON(w, r, &payload); err != nil {
		http.Error(w, "error: invalid credentials", http.StatusInternalServerError)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)

	user := &store.User{
		Username: payload.Username,
		Password: string(hash),
	}

	if err := app.store.Users.Create(user); err != nil {
		log.Println(err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, "message: user successfully created!")
}

func (app *application) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := app.store.Users.GetByUsername(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "error: invalid credentials", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, fmt.Sprintf("message: %s successfully logged in!", user.Username))
}

func (app *application) HandleTest(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	writeJSON(w, http.StatusOK, cookies)
}
