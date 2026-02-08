package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/catouberos/transit-watcher/internal/aggregator"
	"github.com/catouberos/transit-watcher/internal/config"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/catouberos/transit-watcher/providers/multigo"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := config.LoadConfig(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "cannot load application config", "error", err)
		os.Exit(1)
	}

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "cannot setup opentelemetry stack", "error", err)
		os.Exit(1)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	client := &http.Client{
		Timeout: 2 * time.Minute,
	}

	// init kafka
	kafka, err := kgo.NewClient(
		kgo.SeedBrokers(config.Kafka.Seeds...),
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot init kafka", "error", err)
		os.Exit(1)
	}
	defer kafka.Close()

	agg := aggregator.NewAggregator(
		kafka,
		gobus.NewClient(client),
		multigo.NewClient(client),
	)
	agg.Aggregate(ctx)

	<-ctx.Done()
	slog.Info("attempt to gracefully shutdown")
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
