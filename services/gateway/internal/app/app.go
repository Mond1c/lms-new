package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Mond1c/lms/gen/go/lms/v1/lmsv1connect"
	"github.com/Mond1c/lms/pkg/obs"
	"github.com/Mond1c/lms/pkg/pg"
	"github.com/Mond1c/lms/services/gateway/internal/config"
	"github.com/Mond1c/lms/services/gateway/internal/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	cfg     *config.Config
	pool    *pgxpool.Pool
	srv     *http.Server
	ready   atomic.Bool
	shutFns []func(context.Context) error
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	a := &App{cfg: cfg}

	obsShutdown, err := obs.Init(ctx, obs.Config{
		ServiceName:  cfg.ServiceName,
		ServiceVer:   cfg.ServiceVer,
		OTLPEndpoint: cfg.OTLPEndpoint,
		LogLevel:     cfg.LogLevel,
		LogFormat:    cfg.LogFormat,
	})
	if err != nil {
		return nil, fmt.Errorf("obs: %w", err)
	}
	a.shutFns = append(a.shutFns, obsShutdown)

	pool, err := pg.Open(ctx, pg.Config{DSN: cfg.DatabaseURL, MaxConnections: 10})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	a.pool = pool
	a.shutFns = append(a.shutFns, func(context.Context) error { pool.Close(); return nil })
	slog.Info("database connected")

	mux := http.NewServeMux()
	path, h := lmsv1connect.NewGatewayServiceHandler(handler.New())
	mux.Handle(path, h)
	mux.HandleFunc("/healthz", a.healthz)
	mux.HandleFunc("/readyz", a.readyz)

	a.srv = &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	slog.Info("starting http", "addr", a.cfg.HTTPAddr)
	a.ready.Store(true)

	errCh := make(chan error, 1)
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		return a.shutdown()
	case err := <-errCh:
		a.ready.Store(false)
		return fmt.Errorf("server: %w", err)
	}
}

func (a *App) shutdown() error {
	a.ready.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.srv.Shutdown(ctx); err != nil {
		slog.Error("http shutdown", "err", err)
	}
	for _, fn := range a.shutFns {
		if err := fn(ctx); err != nil {
			slog.Error("shutdown fn", "err", err)
		}
	}
	return nil
}

func (a *App) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (a *App) readyz(w http.ResponseWriter, r *http.Request) {
	if !a.ready.Load() {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := a.pool.Ping(ctx); err != nil {
		http.Error(w, "db not reachable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}
