package tasks

import (
	"context"
	"encoding/json"

	"codeberg.org/transit-radar/transit-watcher/internal/aggregator"
	"codeberg.org/transit-radar/transit-watcher/pkg/otelhelper"
	"github.com/hibiken/asynq"
)

type aggregateGoBusHandler struct {
	aggregator *aggregator.Aggregator
}

func NewAggregateGoBusHandler(aggregator *aggregator.Aggregator) asynq.Handler {
	return &aggregateGoBusHandler{aggregator}
}

func (h *aggregateGoBusHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return h.aggregator.AggregateGoBus(ctx)
}

type aggregateMultiGoGeolocationHandler struct {
	aggregator *aggregator.Aggregator
}

func NewAggregateMultiGoGeolocationHandler(aggregator *aggregator.Aggregator) asynq.Handler {
	return &aggregateMultiGoGeolocationHandler{aggregator}
}

func (h *aggregateMultiGoGeolocationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var params AggregateMultiGoGeolocationParams
	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	ctx = otelhelper.ContextFromHeader(ctx, t.Headers())

	return h.aggregator.AggregateMultiGoGeolocation(ctx, params.RouteID, params.VariantID, params.Direction)
}

type aggregateTTGTGeolocationHandler struct {
	aggregator *aggregator.Aggregator
}

func NewAggregateTTGTGeolocationHandler(aggregator *aggregator.Aggregator) asynq.Handler {
	return &aggregateTTGTGeolocationHandler{aggregator}
}

func (h *aggregateTTGTGeolocationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var params AggregateTTGTGeolocationParams
	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	ctx = otelhelper.ContextFromHeader(ctx, t.Headers())

	return h.aggregator.AggregateTTGTGeolocation(ctx, params.RouteID, params.VariantID)
}
