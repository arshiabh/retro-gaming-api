package main

import (
	"net/http"

	"github.com/arshiabh/retro-gaming-api/internal/store"
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

	user, err := app.service.UserService.CreateUser(payload.Username, payload.Password)
	if err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
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

	user, err := app.service.UserService.LoginUser(payload.Username, payload.Password)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := app.auth.GenerateToken(user.ID)
	if err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"token": token})
}

func (app *application) HandleTest(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	writeJSON(w, http.StatusOK, cookies)
}
