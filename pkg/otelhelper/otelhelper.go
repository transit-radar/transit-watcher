package otelhelper

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type otelState struct {
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

func Init(ctx context.Context, name string) error {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
	)

	setupState := otelState{}

	// setup tracer
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	setupState.tracer = tracerProvider

	// setup meter
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
	)
	otel.SetMeterProvider(meterProvider)
	setupState.meter = meterProvider

	// setup logger
	logExporter, err := autoexport.NewLogExporter(ctx)
	if err != nil {
		return err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	global.SetLoggerProvider(loggerProvider)
	setupState.logger = loggerProvider

	slog.SetDefault(otelslog.NewLogger("", otelslog.WithLoggerProvider(loggerProvider)))

	// setup propagator.
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	state.Store(setupState)

	return nil
}

func Shutdown(ctx context.Context) error {
	v := state.Load()
	switch s := v.(type) {
	case nil:
		return errors.New("opentelemetry hasn't been initialised")
	case otelState:
		return shutdown(ctx, s)
	default:
		return errors.New("invalid state")
	}
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
