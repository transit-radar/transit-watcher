package tasks

import (
	"context"
	"log/slog"
	"maps"
	"time"

	"codeberg.org/transit-radar/transit-watcher/internal/clients"
	"codeberg.org/transit-radar/transit-watcher/internal/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/store"
	"codeberg.org/transit-radar/transit-watcher/pkg/otelhelper"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type scheduleMultiGoGeolocationHandler struct {
	store       store.Store
	asynqClient *asynq.Client
}

func NewScheduleMultiGoGeolocationHandler(clients *clients.Clients) asynq.Handler {
	return &scheduleMultiGoGeolocationHandler{
		store:       clients.Store,
		asynqClient: clients.Asynq,
	}
}

func (p *scheduleMultiGoGeolocationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	ctx = otelhelper.ContextFromHeader(ctx, t.Headers())

	slog.InfoContext(ctx, "processing multigo scheduler")

	routes, err := p.store.Members(ctx, v1beta1.CacheKey("routes", "bus"))
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "scheduling geolocation tasks for routes", "routes", routes)

	for _, route := range routes {
		variants, err := p.store.Members(ctx, v1beta1.CacheKey("route", route, "variants"))
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "scheduling geolocation tasks for variants", "route", route, "variants", variants)

		for i, variant := range variants {
			slog.InfoContext(ctx, "scheduling multigo geolocation tasks", "route", route, "direction", i)

			task, err := NewAggregateMultiGoGeolocationTask(AggregateMultiGoGeolocationParams{
				RouteID:   route,
				VariantID: variant,
				Direction: i,
			})
			if err != nil {
				return err
			}

			ctx, _ := otel.GetTracerProvider().Tracer("transit-watcher").Start(ctx, "task/multigo-geolocation")
			carrier := propagation.MapCarrier{}
			otel.GetTextMapPropagator().Inject(ctx, &carrier)
			maps.Copy(task.Headers(), carrier)

			_, err = p.asynqClient.EnqueueContext(ctx, task, asynq.MaxRetry(0), asynq.Unique(30*time.Second))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type scheduleTTGTGeolocationHandler struct {
	store       store.Store
	asynqClient *asynq.Client
}

func NewScheduleTTGTGeolocationHandler(clients *clients.Clients) asynq.Handler {
	return &scheduleTTGTGeolocationHandler{
		store:       clients.Store,
		asynqClient: clients.Asynq,
	}
}

func (p *scheduleTTGTGeolocationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	ctx = otelhelper.ContextFromHeader(ctx, t.Headers())

	slog.InfoContext(ctx, "processing ttgt scheduler")

	routes, err := p.store.Members(ctx, v1beta1.CacheKey("routes", "metro"))
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "scheduling geolocation tasks for routes", "routes", routes)

	for _, route := range routes {
		variants, err := p.store.Members(ctx, v1beta1.CacheKey("route", route, "variants"))
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "scheduling geolocation tasks for variants", "route", route, "variants", variants)

		for _, variant := range variants {
			slog.InfoContext(ctx, "scheduling ttgt geolocation tasks", "route", route, "variant", variant)

			task, err := NewAggregateTTGTGeolocationTask(AggregateTTGTGeolocationParams{
				RouteID:   route,
				VariantID: variant,
			})
			if err != nil {
				return err
			}

			ctx, _ := otel.GetTracerProvider().Tracer("transit-watcher").Start(ctx, "task/ttgt-geolocation")
			carrier := propagation.MapCarrier{}
			otel.GetTextMapPropagator().Inject(ctx, &carrier)
			maps.Copy(task.Headers(), carrier)

			_, err = p.asynqClient.EnqueueContext(ctx, task, asynq.MaxRetry(0), asynq.Unique(1*time.Second))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
