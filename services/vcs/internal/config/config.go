package config

import (
	"fmt"

	"github.com/Mond1c/lms/pkg/obs"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr string `env:"VCS_HTTP_ADDR" envDefault:":8084"`

	DatabaseURL string `env:"VCS_DATABASE_URL,required"`

	ServiceName  string        `env:"VCS_SERVICE_NAME" envDefault:"vcs-svc"`
	ServiceVer   string        `env:"VCS_SERVICE_VERSION" envDefault:"dev"`
	OTLPEndpoint string        `env:"OTLP_ENDPOINT" envDefault:""`
	LogLevel     obs.LogLevel  `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat    obs.LogFormat `env:"LOG_FORMAT" envDefault:"json"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	return &c, nil
}
