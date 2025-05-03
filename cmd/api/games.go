package main

import "net/http"

type gamePayload struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=5"`
}

func (app *application) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	var payload gamePayload
	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var userIDKey contextKey = "userID"
	userID := r.Context().Value(userIDKey).(int64)

	game, err := app.service.GameService.CreateGame(payload.Name, payload.Description, userID)
	if err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": game.ID, "name": game.Name, "description": game.Description})
}
