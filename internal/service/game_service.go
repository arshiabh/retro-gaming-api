package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type GameService struct {
	store *store.Storage
	kafka *kafka.KafkaService
}

func NewGameService(store *store.Storage, kafka *kafka.KafkaService) *GameService {
	return &GameService{
		store: store,
		kafka: kafka,
	}
}

func (s *GameService) CreateGame(name, description string, userID int64) (*store.Game, error) {
	user, err := s.store.Users.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin {
		return nil, ErrForbidden
	}

	game := &store.Game{
		Name:        name,
		Description: description,
	}

	game, err = s.store.Games.Create(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

