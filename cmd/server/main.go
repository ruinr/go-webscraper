package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rnidev/go-webscraper/pkg/cmd"
)

func main() {
	redisPort := flag.String("redisport", "", "port redis listens to")
	gRPCPort := flag.String("grpcport", "", "port grpc listens to")
	gatewayPort := flag.String("gatewayport", "", "port gateway listens to")
	flag.Parse()

	var cfg cmd.Config

	cfg.RedisPort = *redisPort
	cfg.GRPCPort = *gRPCPort
	cfg.RESTPort = *gatewayPort

	if err := cmd.StartServer(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
