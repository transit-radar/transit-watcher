package tasks

import (
	"context"
	"log/slog"

	"codeberg.org/transit-radar/transit-watcher/internal/clients"
	"codeberg.org/transit-radar/transit-watcher/internal/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/store"
	"github.com/hibiken/asynq"
)

type scheduleMultiGoGeolocationHandler struct {
	store       store.Store
	asynqClient *asynq.Client
}

func NewScheduleMultiGoProcessorHandler(clients *clients.Clients) asynq.Handler {
	return &scheduleMultiGoGeolocationHandler{
		store:       clients.Store,
		asynqClient: clients.Asynq,
	}
}

func (p *scheduleMultiGoGeolocationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	slog.InfoContext(ctx, "processing scheduler")
	routes, err := p.store.Members(ctx, v1beta1.CacheKey("routes"))
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

			_, err = p.asynqClient.EnqueueContext(ctx, task)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
