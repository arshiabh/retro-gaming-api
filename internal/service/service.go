package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type Service struct {
	UserService  *UserService
	GameService  *GameService
	ScoreService *ScoreService
}

func NewService(store *store.Storage, kafka kafka.KafkaProducer) *Service {
	return &Service{
		UserService:  NewUserService(store, kafka),
		GameService:  NewGameService(store, kafka),
		ScoreService: NewScoreService(store, kafka),
	}
}
