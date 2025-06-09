package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/arshiabh/retro-gaming-api/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore interface {
	Get(context.Context, int64) (*store.User, error)
	Set(context.Context, *store.User) error
}

type RedisUserStore struct {
	db *redis.Client
}

func NewUserStore(db *redis.Client) *RedisUserStore {
	return &RedisUserStore{
		db: db,
	}
}

func (s *RedisUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cachekey := fmt.Sprintf("user-%d", userID)
	val, err := s.db.Get(ctx, cachekey).Result()
	if err == redis.Nil {
		log.Println("user cache miss")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if val != "" {
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *RedisUserStore) Set(ctx context.Context, user *store.User) error {
	cachekey := fmt.Sprintf("user-%d", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := s.db.SetEx(ctx, cachekey, json, time.Minute*5).Err(); err != nil {
		return err
	}

	return nil
}
