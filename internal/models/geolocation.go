package models

import (
	"time"
)

type Geolocation struct {
	Degree    float32   `redis:"-"`
	Location  LatLng    `redis:"-"`
	Speed     float32   `redis:"-"`
	VehicleID Identity  `redis:"-"`
	RouteID   Identity  `redis:"-"`
	VariantID Identity  `redis:"-"`
	Timestamp time.Time `redis:"timestamp"`
	Hash      uint64    `redis:"hash"`
}

type LatLng struct {
	Latitude  float64
	Longitude float64
}
