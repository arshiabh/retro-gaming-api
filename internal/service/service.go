package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/module"
)

type Service struct {
	UserService  *UserService
	GameService  *GameService
	ScoreService *ScoreService
}

func NewService(deps module.Dependencies) *Service {
	return &Service{
		UserService:  NewUserService(deps),
		GameService:  NewGameService(deps),
		ScoreService: NewScoreService(deps),
	}
}
