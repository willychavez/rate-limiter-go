package core

import (
	"context"
	"time"
)

type RateLimiter struct {
	store            PersistenceInterface
	defaultLimit     int
	defaultBlockTime time.Duration
}

func NewRateLimiter(store PersistenceInterface, defaultLimit int, defaultBlockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		store:            store,
		defaultLimit:     defaultLimit,
		defaultBlockTime: defaultBlockTime,
	}
}

type RateLimiterInterface interface {
	Limit(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, time.Duration)
	GetLimits(ctx context.Context, token string, isToken bool) (int, time.Duration, error)
}

func (r *RateLimiter) Limit(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, time.Duration) {
	count, err := r.store.Increment(ctx, key)
	if err != nil {
		return false, 0
	}

	if count == 1 {
		_ = r.store.Expire(ctx, key, blockTime)
	}

	if count > limit {
		retryAfter := blockTime
		return true, retryAfter
	}

	return false, 0
}

func (r *RateLimiter) GetLimits(ctx context.Context, token string, isToken bool) (int, time.Duration, error) {
	if isToken {
		return r.store.GetTokenLimit(ctx, token)
	}

	return r.defaultLimit, r.defaultBlockTime, nil
}
