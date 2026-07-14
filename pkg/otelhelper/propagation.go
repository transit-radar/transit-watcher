package otelhelper

import (
	"context"
	"maps"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func ContextFromHeader(ctx context.Context, header map[string]string) context.Context {
	carrier := propagation.MapCarrier{}
	maps.Copy(carrier, header)
	return otel.GetTextMapPropagator().Extract(ctx, &carrier)
}
