package redisrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/willychavez/rate-limiter-go/internal/core"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) core.PersistenceInterface {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Increment(ctx context.Context, key string) (int, error) {
	count, err := r.client.Incr(ctx, key).Result()
	return int(count), err
}

func (r *RedisRepository) Expire(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Expire(ctx, key, duration).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (int, error) {
	value, err := r.client.Get(ctx, key).Int()
	return value, err
}

func (r *RedisRepository) Set(ctx context.Context, key string, value int, duration time.Duration) error {
	return r.client.Set(ctx, key, value, duration).Err()
}

func (r *RedisRepository) GetTokenLimit(ctx context.Context, token string) (int, time.Duration, error) {
	limitKey := fmt.Sprintf("token:limit:%s", token)
	blockTimeKey := fmt.Sprintf("token:block_time:%s", token)

	limit, err := r.client.Get(ctx, limitKey).Int()
	if err != nil {
		return 0, 0, fmt.Errorf("token limit not found: %v", err)
	}

	blockTime, err := r.client.Get(ctx, blockTimeKey).Int()
	if err != nil {
		return 0, 0, fmt.Errorf("token block time not found: %v", err)
	}

	return limit, time.Duration(blockTime) * time.Second, nil
}
