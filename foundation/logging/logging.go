package logging

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	DevLogEnv  string = "dev"
	ProdLogEnv string = "prod"
)

func New(service string, env string) (*zap.Logger, error) {
	return loggerFromEnv(service, env)
}

func loggerFromEnv(service string, env string) (*zap.Logger, error) {
	initialFields := map[string]any{
		"service": service,
	}
	switch env {
	case DevLogEnv:
		config := zap.NewDevelopmentConfig()
		config.InitialFields = initialFields
		return zap.NewDevelopment()
	case ProdLogEnv:
		config := zap.NewProductionConfig()
		config.InitialFields = initialFields
		return zap.NewProduction()
	default:
		return nil, fmt.Errorf("unknown environment %s", env)
	}
}
