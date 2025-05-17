package service

import (
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type ScoreService struct {
	kafka kafka.KafkaProducer
	store *store.Storage
}

func NewScoreService(store *store.Storage, kafka kafka.KafkaProducer) *ScoreService {
	return &ScoreService{
		store: store,
		kafka: kafka,
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

func (s *ScoreService) GetLeaderBoard(gameID int64) ([]store.LeaderBoard, error) {
	users, err := s.store.Scores.GetTopTen(gameID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
