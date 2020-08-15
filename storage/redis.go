package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/vsmoraes/open-ledger/ledger/account"
)

const KeyFormat = "lock-account-%s"
const TTL = 10

type RedisClient struct {
	redis     redis.Client
	keyformat string
	ttlms     time.Duration
}

func NewRedisClient(redis redis.Client) *RedisClient {
	return &RedisClient{
		redis: redis,
	}
}

func (cli *RedisClient) key(a *account.Account) string {
	if cli.keyformat != "" {
		return fmt.Sprintf(cli.keyformat, a.ID)
	}

	return fmt.Sprintf(KeyFormat, a.ID)
}

func (cli *RedisClient) ttl() time.Duration {
	if cli.ttlms != 0 {
		return cli.ttlms * time.Second
	}

	return TTL * time.Second
}

func (cli *RedisClient) Lock(ctx context.Context, a *account.Account, key string) error {
	if _, err := cli.redis.SetNX(ctx, cli.key(a), key, cli.ttl()).Result(); err != nil {
		return err
	}

	return nil
}

func (cli *RedisClient) Unlock(ctx context.Context, a *account.Account, key string) error {
	lockedID, err := cli.redis.Get(ctx, cli.key(a)).Result()
	if err != nil {
		return err
	}

	if lockedID != key {
		return errors.New("unable to unlock. This lock belongs to someone else")
	}

	if _, err := cli.redis.Del(ctx, cli.key(a)).Result(); err != nil {
		return err
	}

	return nil
}
