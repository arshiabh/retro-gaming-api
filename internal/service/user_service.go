package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/auth"
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/arshiabh/retro-gaming-api/internal/utils"
)

type UserService struct {
	store *store.Storage
	kafka *kafka.Client
	auth  auth.Authenticator
}

func NewUserService(deps module.Dependencies) *UserService {
	return &UserService{
		store: deps.Store,
		kafka: deps.Kafka,
		auth:  deps.Auth,
	}
}

func (s *UserService) CreateUser(user *store.User) {
	s.store.Users.Create(user)
}

func (s *UserService) LoginUser(username string, password string) (string, error) {
	user, err := s.store.Users.GetByUsername(username)
	if err != nil {
		return "", err
	}

	if err := utils.CheckPasswordHash(user.Password, password); err != nil {
		return "", err
	}

	return s.auth.GenerateToken(user.ID)
}
