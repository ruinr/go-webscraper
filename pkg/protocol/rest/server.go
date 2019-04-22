package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	v1 "github.com/rnidev/go-webscraper/pkg/api/v1"
)

// StartRESTGateWay runs REST gateway for gRPC server
func StartRESTGateWay(ctx context.Context, gRPCPort string, restPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	//ToDo: need to add middleware for authendication between REST and gRPC
	opts := []grpc.DialOption{grpc.WithInsecure()}
	//register gRPC endpoint
	if err := v1.RegisterWebScraperHandlerFromEndpoint(ctx, mux, "localhost:"+gRPCPort, opts); err != nil {
		log.Fatalf("failed to start scraper REST gateway: %v", err)
	}

	server := &http.Server{
		Addr:    ":" + restPort,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down scraper REST gateway...")
			server.Shutdown(ctx)
			<-ctx.Done()
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = server.Shutdown(ctx)
	}()

	log.Println("starting REST gateway...")
	return server.ListenAndServe()
}
