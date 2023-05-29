package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "HEARTBEAT"

type Config struct {
	Service               string        `split_words:"true" default:"HEARTBEAT"`
	Environment           string        `split_words:"true" default:"dev"`
	Host                  string        `split_words:"true" default:":8080"`
	ServerShutdownTimeout time.Duration `split_words:"true" default:"10s"`
	ServerReadTimeout     time.Duration `split_words:"true" default:"250ms"`
	ServerWriteTimeout    time.Duration `split_words:"true" default:"250ms"`
	RedisAddress          string        `split_words:"true" default:"localhost:6379"`
	RedisUsername         string        `split_words:"true" default:""`
	RedisPassword         string        `split_words:"true" default:""`
	RedisReadTimeout      time.Duration `split_words:"true" default:"250ms"`
	RedisWriteTimeout     time.Duration `split_words:"true" default:"250ms"`
	HeartbeatExpiry       time.Duration `split_words:"true" default:"1h"`
	TopicURL              string        `split_words:"true" default:"heartbeat-events"`
	KafkaAddresses        []string      `split_words:"true" default:"localhost:9092"`
}

func loadConfig() (Config, error) {
	var config Config
	err := envconfig.Process(ConfigPrefix, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config from env vars. %w", err)
	}

	return config, nil
}
