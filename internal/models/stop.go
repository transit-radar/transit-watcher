package models

type Stop struct {
	ID       Identity
	Code     string
	Name     string
	Location LatLng

	Hash uint64
}
