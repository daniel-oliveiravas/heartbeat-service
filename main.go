package main

import (
	"log"

	"github.com/daniel-oliveiravas/heartbeat-service/foundation/logger"
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

	log, err := logger.New(cfg.Service, cfg.Environment)
	if err != nil {
		return err
	}

	log.Info("starting service", zap.String("service", cfg.Service))

	return nil
}
