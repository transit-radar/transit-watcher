package gobus

import "encoding/json"

type Route struct {
	Id       json.Number    `json:"_id"`
	Number   string         `json:"number"`
	Name     string         `json:"name"`
	Info     RouteInfo      `json:"info"`
	Variants []RouteVariant `json:"vars"`
}

type RouteInfo struct {
	Id                  json.Number `json:"_id"`
	InboundDescription  string      `json:"inBoundDescription"`
	OutboundDescription string      `json:"outBoundDescription"`
	OperationTime       string      `json:"operationTime"`
	Organization        string      `json:"orgs"`
	Ticketing           string      `json:"tickets"`
	Duration            string      `json:"timeOfTrip"`
	TotalTrip           string      `json:"totalTrip"`
	RouteType           string      `json:"busType"`
}

type RouteVariant struct {
	Id         json.Number        `json:"_id"`
	RouteId    json.Number        `json:"routeId"`
	Name       string             `json:"name"`
	ShortName  string             `json:"shortName"`
	Distance   float32            `json:"distance"`
	StartStop  string             `json:"startStop"`
	EndStop    string             `json:"endStop"`
	IsOutbound bool               `json:"isOutbound"`
	Duration   int32              `json:"runningTime"`
	Stops      []RouteVariantStop `json:"stops"`
}

type RouteVariantWithDescription struct {
	RouteVariant

	Description string `json:"-"`
}

type RouteVariantStop struct {
	Id            json.Number `json:"_id"`
	Code          string      `json:"code"`
	Name          string      `json:"name"`
	Routes        string      `json:"routes"`
	Type          string      `json:"stopType"`
	AddressNumber string      `json:"addressNo"`
	AddressStreet string      `json:"street"`
	AddressWard   string      `json:"ward"`
	AddressZone   string      `json:"zone"`
	Latitude      float32     `json:"lat"`
	Longitude     float32     `json:"lng"`
}

type StopResponse struct {
	// here we dont care about additional geojson features
	Stops []Stop `json:"features"`
}

type Stop struct {
	Geometry StopGeometry `json:"geometry"`
	Property StopProperty `json:"properties"`
}

type StopGeometry struct {
	Coordinates [2]float32 `json:"coordinates"`
}

type StopProperty struct {
	Id       json.Number `json:"stopId"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	TypeName string      `json:"stopType"`
}
