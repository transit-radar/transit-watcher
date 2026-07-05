package models

type Route struct {
	ID          Identity `redis:"-"`
	Number      string   `redis:"number"`
	Name        string   `redis:"name"`
	ShortName   string   `redis:"shortname"`
	Description string   `redis:"description"`
	AgencyID    Identity `redis:"-"`
	Active      bool     `redis:"active"`
	Hash        uint64   `redis:"hash"`
}
