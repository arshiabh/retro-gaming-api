package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type GameService struct {
	store *store.Storage
	kafka *kafka.Client
}

func NewGameService(deps module.Dependencies) *GameService {
	return &GameService{
		store: deps.Store,
		kafka: deps.Kafka,
	}
}

func (s *GameService) CreateGame(name, description string) (*store.Game, error) {
	game := &store.Game{
		Name:        name,
		Description: description,
	}

	game, err := s.store.Games.Create(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}
