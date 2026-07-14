package models

type Variant struct {
	ID     Identity `redis:"-"`
	Number string   `redis:"number"`
	Name   string   `redis:"name"`

	Hash uint64 `redis:"hash"`
}
