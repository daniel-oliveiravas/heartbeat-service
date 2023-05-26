package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "HEARTBEAT"

type Config struct {
	Service     string `default:"HEARTBEAT"`
	Environment string `default:"dev"`
}

func loadConfig() (Config, error) {
	var config Config
	err := envconfig.Process(ConfigPrefix, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config from env vars. %w", err)
	}

	return config, nil
}
