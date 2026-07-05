package models

type Variant struct {
	ID      Identity `redis:"-"`
	Number  string   `redis:"number"`
	Name    string   `redis:"name"`
	RouteID Identity `redis:"-"`

	Hash uint64 `redis:"hash"`
}
