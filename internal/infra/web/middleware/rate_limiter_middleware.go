package middleware

import (
	"net/http"
	"strings"

	"github.com/willychavez/rate-limiter-go/internal/usecases"
)

type RateLimiterMiddleware struct {
	rateLimiterUseCase usecases.RateLimiterUseCaseInterface
}

func NewRateLimiterMiddleware(useCase usecases.RateLimiterUseCaseInterface) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{rateLimiterUseCase: useCase}
}

func (m *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apiKey := r.Header.Get("API_KEY")
		key := getClientIP(r)
		isToken := false

		if apiKey != "" {
			key = strings.TrimSpace(apiKey)
			isToken = true
		}

		// log.Printf("key: %s, isToken: %t", key, isToken)

		if limited, retryAfter := m.rateLimiterUseCase.CheckRateLimit(ctx, key, isToken); limited {
			w.Header().Set("Retry-After", retryAfter)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	// Remove port if present
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}
