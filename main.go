package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/daniel-oliveiravas/heartbeat-service/app/api"
	"github.com/daniel-oliveiravas/heartbeat-service/foundation/logging"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal()
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger, err := logging.New(cfg.Service, cfg.Environment)
	if err != nil {
		return err
	}

	// ----------------------------------------------------------------------------------------------
	// Start API service
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	heartbeatRoutes := api.New()

	//TODO: Configure read, write and idle timeouts
	server := http.Server{
		Addr:     cfg.Host,
		Handler:  heartbeatRoutes,
		ErrorLog: log.Default(),
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("starting server", zap.String("service", cfg.Service), zap.String("host", cfg.Host))
		serverErrors <- server.ListenAndServe()
	}()

	// ----------------------------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Info("shutting down service", zap.String("signal", sig.String()))
		defer logger.Info("shutdown complete", zap.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("failed to stop server gracefully. %w", err)
		}
	}

	return nil
}
