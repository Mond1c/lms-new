package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mond1c/lms/services/vcs/internal/app"
	"github.com/Mond1c/lms/services/vcs/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(2)
	}

	a, err := app.New(ctx, cfg)
	if err != nil {
		slog.Error("app init", "err", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		slog.Error("app run", "err", err)
		os.Exit(1)
	}
}
