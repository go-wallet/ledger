package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/vsmoraes/open-ledger/ledger/account"
)

const KeyFormat = "lock-account-%s"
const TTL = 10 * time.Second

type RedisClient struct {
	redis     *redis.Client
	keyFormat string
	ttlMs     time.Duration
}

func NewRedisClient(redis *redis.Client) *RedisClient {
	return &RedisClient{
		redis: redis,
	}
}

func (cli *RedisClient) Lock(ctx context.Context, id account.ID, key string) error {
	log.WithFields(log.Fields{
		"id":  id,
		"key": key,
	}).Debug("Trying to lock")

	cmd, err := cli.redis.SetNX(ctx, cli.key(id), key, cli.ttl()).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"key": key,
			"cmd": cmd,
		}).
			WithError(err).
			Error("Error locking account")

		return err
	}

	if !cmd {
		log.WithFields(log.Fields{
			"id":  id,
			"key": key,
			"cmd": cmd,
		}).Debug("Account already locked")

		return errors.New("this account is locked")
	}

	log.WithFields(log.Fields{
		"id":  id,
		"key": key,
		"cmd": cmd,
	}).Debug("Locking successful")

	return nil
}

func (cli *RedisClient) Unlock(ctx context.Context, id account.ID, key string) error {
	return nil
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

func (cli *RedisClient) key(id account.ID) string {
	if cli.keyFormat != "" {
		return fmt.Sprintf(cli.keyFormat, id)
	}

	return fmt.Sprintf(KeyFormat, id)
}

func (cli *RedisClient) ttl() time.Duration {
	if cli.ttlMs != 0 {
		return cli.ttlMs * time.Second
	}

	return TTL
}
