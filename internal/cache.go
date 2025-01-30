package internal

import "github.com/redis/go-redis/v9"

func NewRedisClient(opt *redis.Options) *redis.Client {
	return redis.NewClient(opt)
}
