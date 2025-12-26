package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"buf.build/gen/go/catou/transit-radar/connectrpc/go/api/v1/apiv1connect"
	"github.com/catouberos/transit-watcher/internal/aggregator"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/catouberos/transit-watcher/providers/multigo"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
)

type Config struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type KafkaConfig struct {
	Seeds         []string
	ConsumerGroup string
	ConsumeTopics []string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	viper.SetDefault("kafka.seeds", []string{"localhost:9092"})
	viper.SetDefault("kafka.consumergroup", "my-group-identifier")
	viper.SetDefault("kafka.consumetopics", []string{"foo"})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		slog.WarnContext(ctx, "cannot load config, using defaults...", "error", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		slog.ErrorContext(ctx, "cannot parse config", "error", err)
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

	routeService := apiv1connect.NewRouteServiceClient(http.DefaultClient, "http://localhost:5001")
	variantService := apiv1connect.NewVariantServiceClient(http.DefaultClient, "http://localhost:5001")
	geolocationService := apiv1connect.NewGeolocationServiceClient(http.DefaultClient, "http://localhost:5001")
	stopService := apiv1connect.NewStopServiceClient(http.DefaultClient, "http://localhost:5001")

	client := &http.Client{
		Timeout: 2 * time.Minute,
	}

	// init kafka
	kafka, err := kgo.NewClient(
		kgo.SeedBrokers(config.Kafka.Seeds...),
		kgo.ConsumerGroup(config.Kafka.ConsumerGroup),
		kgo.ConsumeTopics(config.Kafka.ConsumeTopics...),
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot init kafka", "error", err)
		os.Exit(1)
	}
	defer kafka.Close()

	agg := aggregator.NewAggregator(
		kafka,
		routeService,
		variantService,
		geolocationService,
		stopService,
		gobus.NewClient(client),
		multigo.NewClient(client),
	)
	agg.Aggregate(context.Background())

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
