package aggregator

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/transit-radar/transit-watcher/internal/mapper"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"codeberg.org/transit-radar/transit-watcher/internal/processor"
	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
)

// AggregateGoBus retrieve routes, variants, and variant-stop mapping from GoBus
func (a *Aggregator) AggregateGoBus(ctx context.Context) error {
	routes, err := a.goBus.ListRoutes(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to retrieve GoBus routes", "error", err)
		return err
	}
	slog.InfoContext(ctx, "successfully retrieved GoBus routes", "count", len(routes))

	for _, route := range routes {
		logger := slog.With(slog.Any("routeID", route.Id))

		r, err := a.processRoute(ctx, route)
		if err != nil {
			logger.ErrorContext(ctx, "failed to process route", "error", err)
		}

		for _, variant := range route.Variants {
			logger := logger.With(slog.Any("variantID", variant.Id))

			if err := a.processVariant(ctx, r, variant); err != nil {
				logger.ErrorContext(ctx, "failed to process variant", "error", err)
			}
		}
	}

	return nil
}

// AggregateGoBusStops retrieve stops from GoBus
func (a *Aggregator) AggregateGoBusStops(ctx context.Context) error {
	stops, err := a.goBus.ListStops(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to retrieve GoBus stops", "error", err)
		return err
	}
	slog.InfoContext(ctx, "successfully retrieved GoBus stops", "count", len(stops))

	for _, stop := range stops {
		if err := a.processStop(ctx, stop); err != nil {
			slog.ErrorContext(ctx, "failed to process stop", "stopID", stop.Property.Id, "error", err)
		}
	}

	return nil
}

func (a *Aggregator) processRoute(ctx context.Context, route gobus.Route) (models.Route, error) {
	r, err := mapper.MapGoBusRoute(route)
	if err != nil {
		return models.Route{}, err
	}

	err = a.route.Validate(ctx, r)
	if errors.Is(err, processor.ErrStaleData) {
		return r, nil
	}
	if err != nil {
		return r, err
	}

	err = a.route.Publish(ctx, r)
	if err != nil {
		return r, err
	}

	err = a.route.Memoize(ctx, r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (a *Aggregator) processVariant(ctx context.Context, route models.Route, variant gobus.RouteVariant) error {
	v, err := mapper.MapGoBusVariant(variant)
	if err != nil {
		return err
	}
	v.RouteID = route.ID

	err = a.variant.Validate(ctx, v)
	if errors.Is(err, processor.ErrStaleData) {
		return nil
	}
	if err != nil {
		return err
	}

	err = a.variant.Publish(ctx, v)
	if err != nil {
		return err
	}

	err = a.variant.Memoize(ctx, v)
	if err != nil {
		return err
	}

	return nil
}

func (a *Aggregator) processStop(ctx context.Context, stop gobus.Stop) error {
	_, err := mapper.MapGoBusStop(stop)
	if err != nil {
		return err
	}

	return nil
}
