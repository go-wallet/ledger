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
const TTL = 10 * time.Second

type RedisClient struct {
	redis     *redis.Client
	keyformat string
	ttlms     time.Duration
}

func NewRedisClient(redis *redis.Client) *RedisClient {
	return &RedisClient{
		redis: redis,
	}
}

func (cli *RedisClient) key(id account.ID) string {
	if cli.keyformat != "" {
		return fmt.Sprintf(cli.keyformat, id)
	}

	return fmt.Sprintf(KeyFormat, id)
}

func (cli *RedisClient) ttl() time.Duration {
	if cli.ttlms != 0 {
		return cli.ttlms * time.Second
	}

	return TTL
}

func (cli *RedisClient) Lock(ctx context.Context, id account.ID, key string) error {
	cmd, err := cli.redis.SetNX(ctx, cli.key(id), key, cli.ttl()).Result()
	if err != nil {
		return err
	}

	if !cmd {
		return errors.New("this account is locked")
	}

	return nil
}

func (cli *RedisClient) Unlock(ctx context.Context, id account.ID, key string) error {
	lockedID, err := cli.redis.Get(ctx, cli.key(id)).Result()
	if err != nil {
		return err
	}

	if lockedID != key {
		return errors.New("unable to unlock. This lock belongs to someone else")
	}

	if _, err := cli.redis.Del(ctx, cli.key(id)).Result(); err != nil {
		return err
	}

	return nil
}
