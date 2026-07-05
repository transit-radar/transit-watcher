package mapper

import (
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
	"github.com/gohugoio/hashstructure"
)

func MapGoBusStop(stop gobus.Stop) (models.Stop, error) {
	hash, err := hashstructure.Hash(stop, nil)
	if err != nil {
		return models.Stop{}, err
	}

	var lat, lng float64
	if lat, err = stop.Geometry.Coordinates[0].Float64(); err != nil {
		return models.Stop{}, err
	}
	if lng, err = stop.Geometry.Coordinates[1].Float64(); err != nil {
		return models.Stop{}, err
	}

	return models.Stop{
		ID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      stop.Property.Id.String(),
		},
		Code:     stop.Property.Code,
		Name:     stop.Property.Name,
		Location: models.LatLng{Latitude: lat, Longitude: lng},
		Hash:     hash,
	}, nil
}

func MapGoBusRoute(route gobus.Route) (models.Route, error) {
	hash, err := hashstructure.Hash(route, nil)
	if err != nil {
		return models.Route{}, err
	}

	return models.Route{
		ID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      route.Id.String(),
		},
		Number: route.Number,
		Name:   route.Name,
		AgencyID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      route.Info.Organization,
		},
		Hash: hash,
	}, nil
}

func MapGoBusVariant(variant gobus.RouteVariant) (models.Variant, error) {
	hash, err := hashstructure.Hash(variant, nil)
	if err != nil {
		return models.Variant{}, err
	}

	return models.Variant{
		ID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      variant.Id.String(),
		},
		Number: variant.Name,
		Hash:   hash,
	}, nil
}
