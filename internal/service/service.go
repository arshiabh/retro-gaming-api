package service

import (
	"sync"

	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/arshiabh/retro-gaming-api/internal/store/cache"
)

type Service struct {
	UserService  *UserService
	GameService  *GameService
	ScoreService *ScoreService
}

func NewService(store *store.Storage, kafka *kafka.KafkaService, rdb *cache.Storage, wg *sync.WaitGroup) *Service {
	return &Service{
		UserService:  NewUserService(store, kafka, rdb, wg),
		GameService:  NewGameService(store, kafka),
		ScoreService: NewScoreService(store, kafka),
	}
}
