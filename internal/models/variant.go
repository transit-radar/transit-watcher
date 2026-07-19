package models

type Variant struct {
	ID        Identity  `redis:"-"`
	Number    string    `redis:"number"`
	Name      string    `redis:"name"`
	Direction Direction `redis:"-"`

	StopIDs []Identity `redis:"-"`

	Hash uint64 `redis:"hash"`
}
