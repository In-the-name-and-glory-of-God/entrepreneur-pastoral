package storage

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStorage interface {
	BuildKey(prefix CachePrefix, data ...string) string
	Get(ctx context.Context, key string, dest any) error
	GetAndDel(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Scan(ctx context.Context, match string) ([]string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type CachePrefix uint8

const (
	CACHE_PREFIX_SESSION CachePrefix = iota
)

func (p CachePrefix) String() string {
	switch p {
	case CACHE_PREFIX_SESSION:
		return "session"
	default:
		return ""
	}
}

type Cache struct {
	client *redis.Client
}

func NewCacheStorage(client *redis.Client) CacheStorage {
	return &Cache{
		client: client,
	}
}

func (c Cache) BuildKey(prefix CachePrefix, data ...string) string {
	return prefix.String() + ":" + strings.Join(data, ":")
}

func (c Cache) Get(ctx context.Context, key string, dest any) error {
	if err := c.client.HGetAll(ctx, key).Scan(dest); err != nil {
		return err
	}

	return nil
}

func (c Cache) GetAndDel(ctx context.Context, key string, dest any) error {
	if err := c.client.HGetAll(ctx, key).Scan(dest); err != nil {
		return err
	}

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (c Cache) Set(ctx context.Context, key string, val any, expire time.Duration) error {
	if err := c.client.HSet(ctx, key, val).Err(); err != nil {
		return err
	}

	if expire > 0 {
		if err := c.client.Expire(ctx, key, expire).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (c Cache) Del(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (c Cache) Scan(ctx context.Context, match string) ([]string, error) {
	var keys []string
	var cursor uint64

	for {
		var batch []string
		var err error
		batch, cursor, err = c.client.Scan(ctx, cursor, match, 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (c Cache) Exists(ctx context.Context, key string) (bool, error) {
	res, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}
