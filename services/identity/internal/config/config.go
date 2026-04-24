package config

import (
	"fmt"

	"github.com/Mond1c/lms/pkg/obs"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr     string        `env:"IDENTITY_HTTP_ADDR" envDefault:":8081"`
	DatabaseURL  string        `env:"IDENTITY_DATABASE_URL,required"`
	ServiceName  string        `env:"IDENTITY_SERVICE_NAME" envDefault:"identity-svc"`
	ServiceVer   string        `env:"IDENTITY_SERVICE_VER" envDefault:"dev"`
	OTLPEndpoint string        `env:"OTLP_ENDPOINT" envDefault:""`
	LogLevel     obs.LogLevel  `env:"IDENTITY_LOG_LEVEL" envDefault:"info"`
	LogFormat    obs.LogFormat `env:"IDENTITY_LOG_FORMAT" envDefault:"json"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	return &c, nil
}
