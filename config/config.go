package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App     App     `envPrefix:"APP_"`
		HTTP    HTTP    `envPrefix:"HTTP_"`
		Log     Log     `envPrefix:"LOG_"`
		PG      PG      `envPrefix:"PG_"`
		GRPC    GRPC    `envPrefix:"GRPC_"`
		RMQ     RMQ     `envPrefix:"RMQ_"`
		Metrics Metrics `envPrefix:"METRICS_"`
		Swagger Swagger `envPrefix:"SWAGGER_"`
		Redis   Redis   `envPrefix:"REDIS_"`
		S3      S3      `envPrefix:"S3_"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}
	return cfg, nil
}
