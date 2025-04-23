package storage

import (
	"context"
	"time"
)

type Storage interface {
	Increment(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value int64, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
