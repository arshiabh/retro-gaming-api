package service

import (
	"github.com/arshiabh/retro-gaming-api/internal/module"
	"github.com/arshiabh/retro-gaming-api/internal/store"
)

type UserService struct {
	store *store.Storage
}

func NewUserService(deps module.Dependencies) *UserService {
	return &UserService{
		store: deps.Store,
	}
}

func (s *UserService) CreateUser(user *store.User) {
	s.store.Users.Create(user)
}
