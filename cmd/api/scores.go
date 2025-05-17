package main

import (
	"net/http"
	"strconv"
)

type scorePayload struct {
	GameID int64 `json:"game_id" validate:"required"`
	Score  int64 `json:"score" validate:"required"`
}

func (app *application) HandleSetScore(w http.ResponseWriter, r *http.Request) {
	var userIDKey contextKey = "userID"
	userID := r.Context().Value(userIDKey).(int64)

	var payload scorePayload

	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	score, err := app.service.ScoreService.SetScore(userID, payload.GameID, payload.Score)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": score.ID, "user_id": score.UserID, "game_id": score.GameID, "score": score.Point, "submitted_at": score.Submitted_at})
}

func (app *application) HandleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	gameIDstr := r.URL.Query().Get("gameID")
	gameID, _ := strconv.Atoi(gameIDstr)
	app.service.ScoreService.GetLeaderBoard(int64(gameID))
	//get top ten from service
}
