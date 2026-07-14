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
	logger := slog.With(slog.Any("routeID", route.Id))

	r, err := mapper.MapGoBusRoute(route)
	if err != nil {
		return models.Route{}, err
	}

	if err := a.enrichRoute(ctx, &r); err != nil {
		logger.WarnContext(ctx, "failed to enrich route with extra metadata", "error", err)
	}

	if err := a.route.Validate(ctx, r); err != nil {
		if errors.Is(err, processor.ErrStaleData) {
			logger.InfoContext(ctx, "retrieved route is stale, skipping...")
			return r, nil
		}

		return r, err
	}

	if err := a.route.Publish(ctx, r); err != nil {
		return r, err
	}
	logger.DebugContext(ctx, "published route payload")

	if err := a.route.Memoize(ctx, r); err != nil {
		return r, err
	}
	logger.DebugContext(ctx, "memoized retrieved route as latest")

	return r, nil
}

func (a *Aggregator) enrichRoute(ctx context.Context, route *models.Route) error {
	return nil

	logger := slog.With(slog.Any("routeID", route.ID.Value))

	r, err := a.ebms.GetRouteByID(ctx, route.ID.Value)
	if err != nil {
		logger.ErrorContext(ctx, "failed to retrieve EBMS route information", "error", err)
		return err
	}

	route.Color = r.Color
	route.Type = r.Type
	route.Organization = r.Orgs
	route.TripDuration = r.TimeOfTrip
	route.Headway = r.Headway
	route.OutboundName = r.OutBoundName
	route.InboundName = r.InBoundName
	route.OutboundDescription = r.OutBoundDescription
	route.InboundDescription = r.InBoundDescription
	route.TotalTrip = r.TotalTrip
	route.Tickets = r.Tickets
	route.OperationTime = models.OperationTime{
		From: r.OperationTime.From,
		To:   r.OperationTime.To,
	}

	return nil
}

func (a *Aggregator) processVariant(ctx context.Context, route models.Route, variant gobus.RouteVariant) error {
	logger := slog.With(slog.Any("routeID", route.ID.Value), slog.Any("variantID", variant.Id))

	v, err := mapper.MapGoBusVariant(variant)
	if err != nil {
		return err
	}

	if err := a.variant.Validate(ctx, route, v); err != nil {
		if errors.Is(err, processor.ErrStaleData) {
			logger.InfoContext(ctx, "retrieved variant is stale, skipping...")
			return nil
		}

		return err
	}

	if err := a.variant.Publish(ctx, route, v); err != nil {
		return err
	}
	logger.DebugContext(ctx, "published variant payload")

	if err := a.variant.Memoize(ctx, route, v); err != nil {
		return err
	}
	logger.DebugContext(ctx, "memoized retrieved variant as latest")

	return nil
}

func (a *Aggregator) processStop(ctx context.Context, stop gobus.Stop) error {
	_, err := mapper.MapGoBusStop(stop)
	if err != nil {
		return err
	}

	return nil
}
