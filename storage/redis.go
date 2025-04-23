package storage

import (
	"context"
	"time"

	"github.com/go-rate-limiter/config"
	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(cfg *config.Config) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisStorage) Get(ctx context.Context, key string) (int64, error) {
	return r.client.Get(ctx, key).Int64()
}

func (r *RedisStorage) Set(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}
