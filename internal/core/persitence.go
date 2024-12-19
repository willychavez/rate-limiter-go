package core

import (
	"context"
	"time"
)

type PersistenceInterface interface {
	Increment(ctx context.Context, key string) (int, error)
	Expire(ctx context.Context, key string, duration time.Duration) error
	Get(ctx context.Context, key string) (int, error)
	Set(ctx context.Context, key string, value int, duration time.Duration) error
	GetTokenLimit(ctx context.Context, token string) (int, time.Duration, error)
}
