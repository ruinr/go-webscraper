package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rnidev/go-webscraper/pkg/cmd"
)

func main() {
	redisHost := flag.String("redishost", "", "host:port redis listens to")
	redisPassword := flag.String("redispassword", "", "password for redis")
	gRPCPort := flag.String("grpcport", "", "port grpc listens to")
	gatewayPort := flag.String("gatewayport", "", "port gateway listens to")
	flag.Parse()

	var cfg cmd.Config

	cfg.RedisHost = *redisHost
	cfg.RedisPassword = *redisPassword
	cfg.GRPCPort = *gRPCPort
	cfg.RESTPort = *gatewayPort

	//Allow flag to override default envrionment variable
	if cfg.RedisHost == "" {
		cfg.RedisHost = os.Getenv("REDIS_HOST")
	}
	if cfg.RedisPassword == "" {
		cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	}
	if cfg.GRPCPort == "" {
		cfg.RedisPassword = os.Getenv("GRPC_PORT")
	}
	if cfg.RESTPort == "" {
		cfg.RESTPort = os.Getenv("PORT")
	}
	if err := cmd.StartServer(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
