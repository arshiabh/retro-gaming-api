package service

import (
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
	score, err := s.store.Scores.Set(gameID, point)
	if err != nil {
		return nil, err
	}
	return score, nil
}
