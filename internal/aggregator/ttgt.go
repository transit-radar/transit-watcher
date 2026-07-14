package aggregator

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/transit-radar/transit-watcher/internal/mapper"
	"codeberg.org/transit-radar/transit-watcher/internal/processor"
	"codeberg.org/transit-radar/transit-watcher/provider/ttgt"
)

func (a *Aggregator) AggregateTTGTGeolocation(ctx context.Context, routeId, variantId string) error {
	logger := slog.With(
		slog.String("routeID", routeId),
		slog.String("variantID", variantId),
	)

	geolocations, err := a.ttgt.ListTransitVehicles(ctx, ttgt.ListTransitVehiclesParams{
		RouteID:   routeId,
		VariantID: variantId,
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to retrieve geolocation from TTGT", "error", err)
		return err
	}
	logger.InfoContext(ctx, "successfully retrieve geolocation from TTGT", "geolocations", geolocations)

	for _, vehicle := range geolocations.Vehicles {
		if err := a.processTTGTGeolocation(ctx, vehicle, routeId, variantId); err != nil {
			logger.ErrorContext(ctx, "error processing geolocation", "error", err)
		}
	}

	return nil
}

func (a *Aggregator) processTTGTGeolocation(ctx context.Context, vehicle ttgt.TransitVehicle, routeId, variantId string) error {
	logger := slog.With(
		slog.String("routeID", routeId),
		slog.String("variantID", variantId),
	)

	g := mapper.MapTTGTGeolocation(vehicle, routeId, variantId)

	if err := a.geolocation.Validate(ctx, g); err != nil {
		if errors.Is(err, processor.ErrStaleData) {
			logger.InfoContext(ctx, "retrieved geolocation is stale, skipping...")
			return nil
		}

		return err
	}

	if err := a.geolocation.Publish(ctx, g); err != nil {
		return err
	}
	logger.DebugContext(ctx, "published geolocation payload")

	if err := a.geolocation.Memoize(ctx, g); err != nil {
		return err
	}
	logger.DebugContext(ctx, "memoized retrieved geolocation as latest")

	return nil
}
