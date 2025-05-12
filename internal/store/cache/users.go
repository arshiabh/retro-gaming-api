package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type UserStore interface {
	Get(context.Context, int64)
}

type RedisUserStore struct {
	db *redis.Client
}

func NewUserStore(db *redis.Client) *RedisUserStore {
	return &RedisUserStore{
		db: db,
	}
}

func (s *RedisUserStore) Get(ctx context.Context, userID int64) {

}
