package factory

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/vsmoraes/open-ledger/internal/config"
	"github.com/vsmoraes/open-ledger/internal/storage"
	"github.com/vsmoraes/open-ledger/ledger/account"
)

func NewLocker() *account.Locker {
	var instances []*account.LockerClient

	for _, conf := range config.Config().Redis.Instances {
		instances = append(instances, account.NewLockerClient(redisClient(conf.Host, conf.Password)))
	}

	return account.NewLocker(instances)
}

func redisClient(host, password string) *storage.RedisClient {
	rc := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
	})
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		panic(err.Error())
	}

	return storage.NewRedisClient(rc)
}
