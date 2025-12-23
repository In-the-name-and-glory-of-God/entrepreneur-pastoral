package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type CacheStorage interface {
	BuildKey(prefix CachePrefix, data ...string) string
	Get(ctx context.Context, key string, dest any) error
	GetAndDel(ctx context.Context, key string, dest any) error
	GetString(ctx context.Context, key string) (string, error)
	GetStringAndDel(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetString(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Scan(ctx context.Context, match string) ([]string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type CachePrefix uint8

const (
	CACHE_PREFIX_SESSION CachePrefix = iota
	CACHE_PREFIX_EMAIL_VERIFICATION
	CACHE_PREFIX_PASSWORD_RESET
	CACHE_PREFIX_REFRESH_TOKEN
)

func (p CachePrefix) String() string {
	switch p {
	case CACHE_PREFIX_SESSION:
		return "session"
	case CACHE_PREFIX_EMAIL_VERIFICATION:
		return "email_verification"
	case CACHE_PREFIX_PASSWORD_RESET:
		return "password_reset"
	case CACHE_PREFIX_REFRESH_TOKEN:
		return "refresh_token"
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
	// Use Lua script to atomically get and delete the key
	// This ensures that if the data is retrieved, it's also deleted in the same operation
	script := redis.NewScript(`
		local value = redis.call('HGETALL', KEYS[1])
		if next(value) ~= nil then
			redis.call('DEL', KEYS[1])
		end
		return value
	`)

	result, err := script.Run(ctx, c.client, []string{key}).Result()
	if err != nil {
		return err
	}

	// Convert the result to a slice and scan into dest
	if resultSlice, ok := result.([]interface{}); ok {
		// Check if the slice is empty (key didn't exist)
		if len(resultSlice) == 0 {
			return nil
		}

		// Use redis.NewSliceResult to create a compatible result for Scan
		cmd := redis.NewSliceResult(resultSlice, nil)
		if err := cmd.Scan(dest); err != nil {
			return err
		}
	}

	return nil
}

func (c Cache) GetString(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrCacheMiss
		}

		return "", err
	}

	return val, nil
}

func (c Cache) GetStringAndDel(ctx context.Context, key string) (string, error) {
	val, err := c.client.GetDel(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrCacheMiss
		}
		return "", err
	}

	return val, nil
}

func (c Cache) Set(ctx context.Context, key string, val any, expire time.Duration) error {
	// Use TxPipelined to execute HSet and Expire atomically
	// This prevents memory leaks if the application crashes between the two operations
	_, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := pipe.HSet(ctx, key, val).Err(); err != nil {
			return err
		}

		if expire > 0 {
			if err := pipe.Expire(ctx, key, expire).Err(); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (c Cache) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
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
