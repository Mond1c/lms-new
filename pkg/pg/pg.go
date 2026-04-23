package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN             string
	MaxConnections  int32
	MinConnections  int32
	HealthCheckTick time.Duration
}

func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pcfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	if cfg.MaxConnections > 0 {
		pcfg.MaxConns = cfg.MaxConnections
	}
	if cfg.MinConnections > 0 {
		pcfg.MinConns = cfg.MinConnections
	}
	if cfg.HealthCheckTick > 0 {
		pcfg.HealthCheckPeriod = cfg.HealthCheckTick
	}

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	const pingTimeout = 5 * time.Second
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}
