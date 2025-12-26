package aggregator

import (
	"context"
	"log/slog"

	apiv1 "buf.build/gen/go/catou/transit-radar/protocolbuffers/go/api/v1"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/cenkalti/backoff/v5"
)

func (a *aggregator) Aggregate(ctx context.Context) error {
	// get routes, variants, and variants-stops
	routes, err := backoff.Retry(
		ctx,
		func() ([]gobus.Route, error) {
			slog.InfoContext(ctx, "getting routes...")
			return a.goBusClient.GetRoutes(ctx)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(5),
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot get routes", "error", err)
		return err
	}
	slog.InfoContext(ctx, "get routes successfully", "count", len(routes))

	// get stops
	stops, err := backoff.Retry(
		ctx,
		func() ([]gobus.Stop, error) {
			slog.InfoContext(ctx, "getting stops...")
			return a.goBusClient.GetStops(ctx)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(5),
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot get stops", "error", err)
		return err
	}
	slog.InfoContext(ctx, "get stops successfully", "count", len(stops))

	// process routes
	for _, route := range routes {
		apiRoute, err := backoff.Retry(
			ctx,
			func() (*apiv1.Route, error) {
				return a.processRoute(ctx, &route)
			},
			backoff.WithBackOff(backoff.NewExponentialBackOff()),
			backoff.WithMaxTries(3),
		)
		if err != nil {
			slog.ErrorContext(ctx, "cannot process route", "error", err)
			continue
		}

		// TODO: mark unused route as inactive

		for _, variant := range route.Variants {
			description := route.Info.InboundDescription
			if variant.IsOutbound {
				description = route.Info.OutboundDescription
			}

			_, err = a.processVariant(ctx, apiRoute.Id, variant, description)
			if err != nil {
				slog.ErrorContext(ctx, "cannot process variant", "error", err)
				continue
			}
		}
	}

	// process stops

	// create relationship between variants and stops

	return nil
}
