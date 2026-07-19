package models

import "time"

type Direction int

const (
	DirectionUnspecified Direction = iota
	DirectionOneDirection
	DirectionOppositeDirection
)

type Route struct {
	ID          Identity  `redis:"-"`
	Number      string    `redis:"number"`
	Name        string    `redis:"name"`
	ShortName   string    `redis:"shortName"`
	Description string    `redis:"description"`
	AgencyID    Identity  `redis:"-"`
	RouteType   RouteType `redis:"-"`
	Active      bool      `redis:"active"`

	// Metadata
	Color               *string        `redis:"-"`
	FareType            *string        `redis:"-"`
	Distance            *float32       `redis:"-"`
	Organization        *string        `redis:"-"`
	TripDuration        *string        `redis:"-"`
	Headway             *string        `redis:"-"`
	OperationTime       *OperationTime `redis:"-"`
	SeatCounts          *[]int         `redis:"-"`
	OutboundName        *string        `redis:"-"`
	InboundName         *string        `redis:"-"`
	OutboundDescription *string        `redis:"-"`
	InboundDescription  *string        `redis:"-"`
	TotalTrip           *string        `redis:"-"`
	Tickets             *string        `redis:"-"`

	Hash uint64 `redis:"hash"`
}

type OperationTime struct {
	From, To time.Time
}

type RouteType int

const (
	RouteTypeUnspecified RouteType = iota
	RouteTypeTram
	RouteTypeMetro
	RouteTypeRail
	RouteTypeBus
	RouteTypeFerry
	RouteTypeCableTram
	RouteTypeCableCar
	RouteTypeFunicular
	RouteTypeTrolleybus
	RouteTypeMonorail
)
