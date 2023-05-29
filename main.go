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

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:         cfg.RedisAddress,
		Username:     cfg.RedisUsername,
		Password:     cfg.RedisPassword,
		ReadTimeout:  cfg.RedisReadTimeout,
		WriteTimeout: cfg.RedisWriteTimeout,
	})
	redisRepository := redis.NewRepository(redisClient, cfg.HeartbeatExpiry)
	publisherConfig := kafka.PublisherConfig{
		TopicURL:       cfg.TopicURL,
		KafkaAddresses: cfg.KafkaAddresses,
	}
	kafkaPublisher, err := kafka.NewPublisher(publisherConfig)
	if err != nil {
		return err
	}
	heartbeatUsecase := heartbeat.NewUsecase(logger, redisRepository, kafkaPublisher)

	apiCfg := api.HandlerConfig{
		HeartbeatUsecase: heartbeatUsecase,
		Logger:           logger,
	}

	heartbeatRoutes, err := api.New(apiCfg)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:         cfg.Host,
		Handler:      heartbeatRoutes,
		ErrorLog:     log.Default(),
		WriteTimeout: cfg.ServerWriteTimeout,
		ReadTimeout:  cfg.ServerReadTimeout,
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
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerShutdownTimeout)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("failed to stop server gracefully. %w", err)
		}
	}

	return nil
}
