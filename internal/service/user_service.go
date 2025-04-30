package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/arshiabh/retro-gaming-api/internal/utils"
)

type UserService struct {
	store *store.Storage
	kafka *kafka.Client
}

func NewUserService(deps module.Dependencies) *UserService {
	return &UserService{
		store: deps.Store,
		kafka: deps.Kafka,
	}
}

func (s *UserService) CreateUser(username, password string) (*store.User, error) {
	hashpassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &store.User{
		Username: username,
		Password: hashpassword,
	}

	user, err = s.store.Users.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) LoginUser(username string, password string) (*store.User, error) {
	user, err := s.store.Users.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := utils.CheckPasswordHash(user.Password, password); err != nil {
		return nil, err
	}

	return user, nil
}
