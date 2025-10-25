package cache

import (
	"context"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/redis/go-redis/v9"
)

func New(cfg config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Protocol: cfg.Protocol,
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
