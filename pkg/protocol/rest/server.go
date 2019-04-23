package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/rnidev/go-webscraper/pkg/api/v1"
	"github.com/rnidev/go-webscraper/pkg/logger"
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
		logger.Log.Fatal("failed to start scraper REST gateway", zap.String("error", err.Error()))
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
			logger.Log.Warn("shutting down scraper REST gateway...")
			server.Shutdown(ctx)
			<-ctx.Done()
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = server.Shutdown(ctx)
	}()

	logger.Log.Info("starting REST gateway", zap.String("port:", restPort))

	return server.ListenAndServe()
}
