package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "HEARTBEAT"

type Config struct {
	Service           string        `default:"HEARTBEAT"`
	Environment       string        `default:"dev"`
	Host              string        `default:":8080"`
	ShutdownTimeout   time.Duration `default:"10s"`
	RedisAddress      string        `default:"localhost:6379"`
	HeartbeatExpiry   time.Duration `default:"30s"`
	HeartbeatTopicURL string        `default:"heartbeat-events"`
	KafkaAddresses    []string      `default:"localhost:9092"`
}

func loadConfig() (Config, error) {
	var config Config
	err := envconfig.Process(ConfigPrefix, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config from env vars. %w", err)
	}

	return config, nil
}
