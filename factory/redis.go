package factory

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/vsmoraes/open-ledger/config"
	"github.com/vsmoraes/open-ledger/storage"
)

func NewLocker() (*storage.RedisClient, *redis.Client) {
	rc := redis.NewClient(&redis.Options{
		Addr:     config.Config().Redis.Host,
		Password: config.Config().Redis.Password,
	})
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		panic(err.Error())
	}

	redisClient := storage.NewRedisClient(rc)

	return redisClient, rc
}
