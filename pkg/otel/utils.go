package otel

import (
	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	UnknownPackageName = "package.unknown"
)

func Tracer() trace.Tracer {
	name := UnknownPackageName
	state, err := loadState()
	if err == nil {
		name = state.packageName
	}

	return otelsdk.GetTracerProvider().Tracer(name)
}

func Meter() metric.Meter {
	name := UnknownPackageName
	state, err := loadState()
	if err == nil {
		name = state.packageName
	}

	return otelsdk.GetMeterProvider().Meter(name)
}

func TextMapPropagator() propagation.TextMapPropagator {
	return otelsdk.GetTextMapPropagator()
}
