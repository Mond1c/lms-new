package obs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelInfo  LogLevel = "info"
)

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

type Config struct {
	ServiceName  string
	ServiceVer   string
	OTLPEndpoint string
	LogLevel     LogLevel
	LogFormat    LogFormat
}

func Init(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if err := setupLogger(cfg); err != nil {
		return nil, err
	}

	if cfg.OTLPEndpoint == "" {
		slog.Info("tracing disabled: OTLP_ENDPOINT empty")
		return func(context.Context) error { return nil }, nil
	}

	const batchTimeout = time.Second * 5

	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(batchTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVer),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to merge resource attributes: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	slog.Info("tracing initialized", "endpoint", cfg.OTLPEndpoint)
	return tp.Shutdown, nil
}

func setupLogger(cfg Config) error {
	level := parseLevel(cfg.LogLevel)
	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if cfg.LogFormat == LogFormatText {
		h = slog.NewTextHandler(os.Stdout, opts)
	} else if cfg.LogFormat == LogFormatJSON {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		return fmt.Errorf("unknown log format: %s", cfg.LogFormat)
	}
	slog.SetDefault(slog.New(h).With(
		slog.String("service", cfg.ServiceName),
	))
	return nil
}

func parseLevel(l LogLevel) slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
