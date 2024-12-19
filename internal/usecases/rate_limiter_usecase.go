package usecases

import (
	"context"

	"github.com/willychavez/rate-limiter-go/internal/core"
)

type RateLimiterUseCase struct {
	rateLimiter core.RateLimiterInterface
}

func NewRateLimiterUseCase(rateLimiter core.RateLimiterInterface) RateLimiterUseCaseInterface {
	return &RateLimiterUseCase{rateLimiter: rateLimiter}
}

type RateLimiterUseCaseInterface interface {
	CheckRateLimit(ctx context.Context, key string, isToken bool) (bool, string)
}

func (u *RateLimiterUseCase) CheckRateLimit(ctx context.Context, key string, isToken bool) (bool, string) {
	limit, blockTime, err := u.rateLimiter.GetLimits(ctx, key, isToken)
	if err != nil {
		return false, ""
	}

	limited, retryAfter := u.rateLimiter.Limit(ctx, key, limit, blockTime)
	if limited {
		return true, retryAfter.String()
	}

	return false, ""
}
