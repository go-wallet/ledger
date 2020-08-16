package factory

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/vsmoraes/open-ledger/storage"
)

func NewRedisClient() (*storage.RedisClient, *redis.Client) {
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		panic(err.Error())
	}

	redisClient := storage.NewRedisClient(rc)

	return redisClient, rc
}
