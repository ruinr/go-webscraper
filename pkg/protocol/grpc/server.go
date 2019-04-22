package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v1 "github.com/rnidev/go-webscraper/pkg/api/v1"
)

// StartgRPCServer runs gRPC service to publish scraper server
func StartgRPCServer(ctx context.Context, srvHandler v1.WebScraperServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register scraper server
	server := grpc.NewServer()
	v1.RegisterWebScraperServer(server, srvHandler)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down scraper gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start scraper gRPC server
	log.Println("starting scraper gRPC server...")
	return server.Serve(listen)
}
