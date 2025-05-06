package service

import (
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type ScoreService struct {
	kafka *kafka.Client
	store *store.Storage
}

func NewScoreService(deps module.Dependencies) *ScoreService {
	return &ScoreService{
		kafka: deps.Kafka,
		store: deps.Store,
	}
}

func (s *ScoreService) SetScore(userID, gameID, point int64) (*store.Score, error) {
	game, err := s.store.Games.GetByGameID(gameID)
	if err != nil {
		return nil, err
	}
	score := &store.Score{
		UserID:       userID,
		GameID:       game.ID,
		Point:        point,
		Submitted_at: time.Now(),
	}

	score, err = s.store.Scores.Set(score)
	if err != nil {
		return nil, err
	}
	return score, nil
}
