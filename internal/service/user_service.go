package service

import (
	"fmt"

	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/arshiabh/retro-gaming-api/internal/utils"
)

type UserService struct {
	store *store.Storage
	kafka kafka.KafkaProducer
}

func NewUserService(store *store.Storage, kafka kafka.KafkaProducer) *UserService {
	return &UserService{
		store: store,
		kafka: kafka,
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

	if err := s.kafka.EnsureTopicExists("user-signup"); err != nil {
		return nil, err
	}
	if err := s.kafka.Produce("user-signup",
		fmt.Appendf(nil, "%d", user.ID),
		fmt.Appendf(nil, `{"event":"user-signup", "user_id":%d, "username":%v }`, user.ID, user.Username)); err != nil {
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
