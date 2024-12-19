package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/willychavez/rate-limiter-go/config/database/redisdb"
	"github.com/willychavez/rate-limiter-go/internal/core"
	"github.com/willychavez/rate-limiter-go/internal/infra/database/redisrepository"
	"github.com/willychavez/rate-limiter-go/internal/infra/web/httpserver"
	"github.com/willychavez/rate-limiter-go/internal/usecases"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	// Graceful shutdown setup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis client
	redisClient := redisdb.NewRedisClient()
	defer redisClient.Close()

	// Initialize Redis repository
	redisRepo := redisrepository.NewRedisRepository(redisClient)

	// Rate limiter configuration
	requestLimit, err := strconv.Atoi(os.Getenv("REQUEST_LIMIT"))
	if err != nil {
		requestLimit = 10 // Default limit
	}
	blockTime, err := strconv.Atoi(os.Getenv("BLOCK_TIME"))
	if err != nil {
		blockTime = 60 // Default block time in seconds
	}

	// Initialize Rate Limiter
	rateLimiter := core.NewRateLimiter(redisRepo, requestLimit, time.Duration(blockTime)*time.Second)
	rateLimiterUseCase := usecases.NewRateLimiterUseCase(rateLimiter)

	// Initialize HTTP Server
	server := httpserver.NewHTTPServer(rateLimiterUseCase)

	// Start server
	log.Println("Starting server on :8080")
	go func() {
		if err := server.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down server...")
	}
}
