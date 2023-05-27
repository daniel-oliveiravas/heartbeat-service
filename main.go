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
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat/integration/kafka"
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat/integration/redis"
	"github.com/daniel-oliveiravas/heartbeat-service/foundation/logging"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
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

	//TODO: Configure redis timeouts and password
	redisClient := goredis.NewClient(&goredis.Options{
		Addr: cfg.RedisAddress,
	})
	redisRepository := redis.NewRepository(redisClient, cfg.HeartbeatExpiry)
	publisherConfig := kafka.PublisherConfig{
		TopicURL:       cfg.HeartbeatTopicURL,
		KafkaAddresses: cfg.KafkaAddresses,
	}
	kafkaPublisher, err := kafka.NewPublisher(publisherConfig)
	if err != nil {
		return err
	}
	heartbeatUsecase := heartbeat.NewUsecase(logger, redisRepository, kafkaPublisher)

	apiCfg := api.HandlerConfig{
		HeartbeatUsecase: heartbeatUsecase,
	}

	heartbeatRoutes, err := api.New(apiCfg)
	if err != nil {
		return err
	}

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
