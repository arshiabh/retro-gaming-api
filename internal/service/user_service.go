package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/events"
	"github.com/arshiabh/retro-gaming-api/internal/kafka"
	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/arshiabh/retro-gaming-api/internal/store/cache"
	"github.com/arshiabh/retro-gaming-api/internal/utils"
)

type UserService struct {
	store *store.Storage
	kafka *kafka.KafkaService
	rdb   *cache.Storage
	wg    *sync.WaitGroup
}

func NewUserService(store *store.Storage, kafka *kafka.KafkaService, rdb *cache.Storage, wg *sync.WaitGroup) *UserService {
	return &UserService{
		store: store,
		kafka: kafka,
		rdb:   rdb,
		wg:    wg,
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

	event := &events.SignedUpEvent{
		EventType: "User Created",
		UserID:    strconv.FormatInt(user.ID, 10),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	marshalEvent, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	kafka.SendAsync(s.wg, "user-signup", fmt.Sprintf("%d", user.ID), marshalEvent, s.kafka)

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
