package aggregator

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/transit-radar/transit-watcher/internal/mapper"
	"codeberg.org/transit-radar/transit-watcher/internal/processor"
	"codeberg.org/transit-radar/transit-watcher/provider/multigo"
)

func (a *Aggregator) AggregateMultiGoGeolocation(ctx context.Context, routeId, variantId string, direction int) error {
	logger := slog.With(
		slog.String("routeID", routeId),
		slog.String("variantID", variantId),
		slog.Int("direction", direction),
	)

	geolocations, err := a.multiGo.ListGeolocations(ctx, multigo.ListGeolocationParams{
		RouteID:   &routeId,
		Direction: &direction,
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to retrieve geolocation from MultiGo", "error", err)
		return err
	}
	logger.InfoContext(ctx, "successfully retrieve geolocation from MultiGo", "geolocations", geolocations)

	for _, geolocation := range geolocations {
		if err := a.processGeolocation(ctx, geolocation, variantId); err != nil {
			logger.ErrorContext(ctx, "error processing geolocation", "error", err)
		}
	}

	return nil
}

func (a *Aggregator) processGeolocation(ctx context.Context, geolocation multigo.Geolocation, variantId string) error {
	g, err := mapper.MapMultiGoGeolocation(geolocation, variantId)
	if err != nil {
		return err
	}

	if err := a.geolocation.Validate(ctx, g); err != nil {
		if !errors.Is(err, processor.ErrStaleData) {
			return err
		}
	}

	if err := a.geolocation.Publish(ctx, g); err != nil {
		return err
	}

	if err := a.geolocation.Memoize(ctx, g); err != nil {
		return err
	}

	return nil
}
