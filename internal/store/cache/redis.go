package cache

import "github.com/redis/go-redis/v9"

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0, //default
		Password: "",
	})

}
