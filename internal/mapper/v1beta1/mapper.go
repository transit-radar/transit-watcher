package v1beta1

import (
	"errors"

	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/mapper"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"github.com/gohugoio/hashstructure"
	latlng "google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapRoute(route models.Route) (*processorv1beta1.Route, error) {
	hash, err := hashstructure.Hash(route, nil)
	if err != nil {
		return nil, err
	}

	id, err := MapIdentity(route.ID)
	if err != nil {
		return nil, err
	}

	routeType, err := mapper.RouteType(route)
	if err != nil {
		return nil, err
	}

	r := &processorv1beta1.Route{
		Id:     id,
		Number: route.Number,
		Hash:   hash,
		Type:   routeType,
	}

	return r, nil
}

func MapVariant(variant models.Variant) (*processorv1beta1.Variant, error) {
	hash, err := hashstructure.Hash(variant, nil)
	if err != nil {
		return nil, err
	}

	id, err := MapIdentity(variant.ID)
	if err != nil {
		return nil, err
	}

	v := &processorv1beta1.Variant{
		Id:   id,
		Name: variant.Name,
		Hash: hash,
	}

	return v, nil
}

func MapGeolocation(geolocation models.Geolocation) (*processorv1beta1.Geolocation, error) {
	hash, err := hashstructure.Hash(geolocation, nil)
	if err != nil {
		return nil, err
	}

	vehicleID, err := MapIdentity(geolocation.VehicleID)
	if err != nil {
		return nil, err
	}

	variantID, err := MapIdentity(geolocation.VariantID)
	if err != nil {
		return nil, err
	}

	routeID, err := MapIdentity(geolocation.RouteID)
	if err != nil {
		return nil, err
	}

	g := &processorv1beta1.Geolocation{
		Degree: float64(geolocation.Degree),
		Location: &latlng.LatLng{
			Latitude:  geolocation.Location.Latitude,
			Longitude: geolocation.Location.Longitude,
		},
		Speed:     float64(geolocation.Speed),
		VehicleId: vehicleID,
		RouteId:   routeID,
		VariantId: variantID,
		Timestamp: timestamppb.New(geolocation.Timestamp),
		Hash:      hash,
	}

	return g, nil
}

func MapIdentity(identity models.Identity) (*processorv1beta1.Identity, error) {
	var identifier processorv1beta1.ExternalIdentifier
	switch identity.Identifier {
	case models.ExternalIdentifierEBMS:
		identifier = processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_EBMS
	case models.ExternalIdentifierMultiGo:
		identifier = processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_MULTIGO
	default:
		return nil, errors.New("unsupported identifier")
	}

	return &processorv1beta1.Identity{
		Identifier: identifier,
		Value:      identity.Value,
	}, nil
}
