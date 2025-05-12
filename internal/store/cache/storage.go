package cache

import "github.com/redis/go-redis/v9"

type Storage struct {
	User UserStore
}

func NewStorage(db *redis.Client) *Storage {
	return &Storage{
		User: NewUserStore(db),
	}
}
