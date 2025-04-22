package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arshiabh/retro-gaming-api/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type userPayload struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
}

func (app *application) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	var users []store.User
	writeJSON(w, http.StatusOK, users)
}

func (app *application) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload userPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	//validate come form json
	if err := validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	writeJSON(w, http.StatusCreated, map[string]any{"id": user.ID, "username": user.Username})
}

func (app *application) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var payload userPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	//validate come form json
	if err := validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.store.Users.Login(payload.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		writeErrJSON(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, fmt.Sprintf("message: %s successfully logged in!", user.Username))
}

func (app *application) HandleTest(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	writeJSON(w, http.StatusOK, cookies)
}
