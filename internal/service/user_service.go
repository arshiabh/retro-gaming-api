package service

import (
	"context"
	"fmt"
	"log"

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

func (s *UserService) Readmessage(ctx context.Context) {
	reader := s.kafka.CreateReader("user-signup-consumer", "user-signup")
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("kafka consumer shutting down!")
			return
		default:
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("error reading message")
			}
			log.Printf("message recevied: %s\n", string(m.Value))
		}
	}
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
