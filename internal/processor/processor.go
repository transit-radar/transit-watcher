package processor

import (
	"context"
	"errors"

	"codeberg.org/transit-radar/transit-watcher/internal/models"
)

var (
	ErrStaleData = errors.New("stale data")
)

type GeolocationProcessor interface {
	// best-effort deduplication
	Validate(context.Context, models.Geolocation) error
	Publish(context.Context, models.Geolocation) error
	Memoize(context.Context, models.Geolocation) error
}

type RouteProcessor interface {
	Validate(context.Context, models.Route) error
	Publish(context.Context, models.Route) error
	Memoize(context.Context, models.Route) error
}

type VariantProcessor interface {
	Validate(context.Context, models.Route, models.Variant) error
	Publish(context.Context, models.Route, models.Variant) error
	Memoize(context.Context, models.Route, models.Variant) error
}
