package v1beta1

import (
	"errors"

	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	radarv1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/radar/v1"
	"codeberg.org/transit-radar/transit-watcher/internal/mapper"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"github.com/gohugoio/hashstructure"
	"github.com/icza/gox/imagex/colorx"
	"google.golang.org/genproto/googleapis/type/color"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	agencyID, err := MapIdentity(route.AgencyID)
	if err != nil {
		return nil, err
	}

	builder := &processorv1beta1.Route_builder{
		Id:          id,
		AgencyId:    agencyID,
		Number:      route.Number,
		Name:        route.Name,
		ShortName:   route.ShortName,
		Description: route.Description,
		Type:        routeType,
		Hash:        hash,
	}

	if route.Color != nil {
		routeColor, err := colorx.ParseHexColor(*route.Color)
		if err != nil {
			return nil, err
		}
		builder.Color = &color.Color{
			Red:   float32(routeColor.R),
			Green: float32(routeColor.G),
			Blue:  float32(routeColor.B),
			Alpha: &wrapperspb.FloatValue{
				Value: float32(routeColor.A),
			},
		}
	}

	return builder.Build(), nil
}

func MapTrip(route models.Route, variant models.Variant) (*processorv1beta1.Trip, error) {
	hash, err := hashstructure.Hash(variant, nil)
	if err != nil {
		return nil, err
	}

	id, err := MapIdentity(variant.ID)
	if err != nil {
		return nil, err
	}

	routeID, err := MapIdentity(route.ID)
	if err != nil {
		return nil, err
	}

	v := processorv1beta1.Trip_builder{
		Id:        id,
		RouteId:   routeID,
		Headsign:  variant.Name,
		ShortName: &variant.Number,
		Direction: new(MapDirection(variant.Direction)),
		Hash:      hash,
	}

	return v.Build(), nil
}

func MapStop(stop models.Stop) (*processorv1beta1.Stop, error) {
	id, err := MapIdentity(stop.ID)
	if err != nil {
		return nil, err
	}

	stopType, err := MapStopType(stop.Type)
	if err != nil {
		return nil, err
	}

	builder := processorv1beta1.Stop_builder{
		Id:   id,
		Code: stop.Code,
		Name: stop.Name,
		Type: stopType,
		Location: &latlng.LatLng{
			Latitude:  stop.Location.Latitude,
			Longitude: stop.Location.Longitude,
		},
	}

	return builder.Build(), nil
}

func MapTripStop(routeID, variantID, stopID models.Identity, orderScore int32) (*processorv1beta1.TripStop, error) {
	pbRouteID, err := MapIdentity(routeID)
	if err != nil {
		return nil, err
	}

	pbTripID, err := MapIdentity(variantID)
	if err != nil {
		return nil, err
	}

	pbStopID, err := MapIdentity(stopID)
	if err != nil {
		return nil, err
	}

	builder := processorv1beta1.TripStop_builder{
		RouteId:    pbRouteID,
		TripId:     pbTripID,
		StopId:     pbStopID,
		OrderScore: orderScore,
	}

	return builder.Build(), nil
}

func MapStopType(stopType models.StopType) (radarv1.StopType, error) {
	switch stopType {
	case models.StopTypeStopPlatform:
		return radarv1.StopType_STOP_TYPE_STOP_PLATFORM, nil
	case models.StopTypeStation:
		return radarv1.StopType_STOP_TYPE_STATION, nil
	}

	return radarv1.StopType_STOP_TYPE_UNSPECIFIED, errors.New("unhandled stop type")
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

func MapDirection(direction models.Direction) processorv1beta1.Direction {
	switch direction {
	case models.DirectionOneDirection:
		return processorv1beta1.Direction_DIRECTION_ONE_DIRECTION
	case models.DirectionOppositeDirection:
		return processorv1beta1.Direction_DIRECTION_OPPOSITE_DIRECTION
	default:
		return processorv1beta1.Direction_DIRECTION_UNSPECIFIED
	}
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
