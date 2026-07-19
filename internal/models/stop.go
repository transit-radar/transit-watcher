package models

type StopType int

const (
	StopTypeUnspecified StopType = iota
	StopTypeStopPlatform
	StopTypeStation
	StopTypeEntraceExit
	StopTypeGenericNode
	StopTypeBoardingArea
)

type Stop struct {
	ID       Identity `redis:"-"`
	Code     string   `redis:"code"`
	Name     string   `redis:"name"`
	Type     StopType `redis:"-"`
	Location LatLng   `redis:"-"`

	Hash uint64 `redis:"hash"`
}
