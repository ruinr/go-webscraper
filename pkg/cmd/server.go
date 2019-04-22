package cmd

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/rnidev/go-webscraper/pkg/protocol/grpc"
	"github.com/rnidev/go-webscraper/pkg/protocol/rest"
	v1 "github.com/rnidev/go-webscraper/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	GRPCPort  string
	RESTPort  string
	RedisPort string
}

// StartServer runs gRPC server and REST gateway
func StartServer(cfg *Config) error {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:         ":" + cfg.RedisPort,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}).WithContext(ctx)

	v1API := v1.NewScraperServer(client)

	// run REST gateway
	go func() {
		_ = rest.StartRESTGateWay(ctx, cfg.GRPCPort, cfg.RESTPort)

	}()
	return grpc.StartgRPCServer(ctx, v1API, cfg.GRPCPort)
}
