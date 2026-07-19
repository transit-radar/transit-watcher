package otel

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type otelState struct {
	packageName string

	// tracer
	tracer *trace.TracerProvider

	// meter
	meter *metric.MeterProvider

	// logger
	logger *log.LoggerProvider
}

var state = defaultState()

func defaultState() *atomic.Value {
	v := &atomic.Value{}
	v.Store(otelState{})
	return v
}

func Init(ctx context.Context, serviceName, packageName string) error {
	setupState := otelState{packageName: packageName}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	// setup tracer
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
		trace.WithResource(res),
	)
	otelsdk.SetTracerProvider(tracerProvider)
	setupState.tracer = tracerProvider

	// setup meter
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(res),
	)
	otelsdk.SetMeterProvider(meterProvider)
	setupState.meter = meterProvider

	// setup logger
	logExporter, err := autoexport.NewLogExporter(ctx)
	if err != nil {
		return err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)
	global.SetLoggerProvider(loggerProvider)
	setupState.logger = loggerProvider

	// setup propagator.
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otelsdk.SetTextMapPropagator(prop)

	slog.SetDefault(otelslog.NewLogger(packageName))

	state.Store(setupState)

	return nil
}

func loadState() (otelState, error) {
	v := state.Load()
	switch s := v.(type) {
	case nil:
		return otelState{}, errors.New("opentelemetry hasn't been initialised")
	case otelState:
		return s, nil
	default:
		return otelState{}, errors.New("invalid state")
	}
}

func Shutdown(ctx context.Context) error {
	state, err := loadState()
	if err != nil {
		return err
	}

	return shutdown(ctx, state)
}

func shutdown(ctx context.Context, s otelState) error {
	if s.tracer != nil {
		slog.DebugContext(ctx, "shutting down opentelemetry tracer provider")
		if err := s.tracer.Shutdown(ctx); err != nil {
			return err
		}
	}

	if s.meter != nil {
		slog.DebugContext(ctx, "shutting down opentelemetry meter provider")
		if err := s.meter.Shutdown(ctx); err != nil {
			return err
		}
	}

	if s.logger != nil {
		slog.DebugContext(ctx, "shutting down opentelemetry logger provider")
		if err := s.logger.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
