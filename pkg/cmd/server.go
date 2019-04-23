package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/rnidev/go-webscraper/pkg/logger"
	"github.com/rnidev/go-webscraper/pkg/protocol/grpc"
	"github.com/rnidev/go-webscraper/pkg/protocol/rest"
	v1 "github.com/rnidev/go-webscraper/pkg/service/v1"
	"go.uber.org/zap"
)

// Config is configuration for Server
type Config struct {
	GRPCPort  string
	RESTPort  string
	RedisHost string
}

// StartServer runs gRPC server and REST gateway
func StartServer(cfg *Config) error {
	ctx := context.Background()
	// start logger with default LogLevel 0
	if err := logger.Init(0); err != nil {
		return fmt.Errorf("failed to start logger: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisHost,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}).WithContext(ctx)
	err := client.Ping().Err()
	if err != nil {
		logger.Log.Warn("Redis server is not available", zap.String("error:", err.Error()))
	}

	v1API := v1.NewScraperServer(client)

	// run REST gateway
	go func() {
		_ = rest.StartRESTGateWay(ctx, cfg.GRPCPort, cfg.RESTPort)

	}()
	return grpc.StartgRPCServer(ctx, v1API, cfg.GRPCPort)
}
