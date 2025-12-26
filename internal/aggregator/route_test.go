package aggregator

import (
	"testing"

	apiv1 "buf.build/gen/go/catou/transit-radar/protocolbuffers/go/api/v1"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/stretchr/testify/assert"
)

func TestRouteType(t *testing.T) {
	tests := []struct {
		name        string
		routeNumber string
		routeType   apiv1.RouteType
	}{{
		name:        "PublicBus",
		routeNumber: "161",
		routeType:   apiv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "PublicBusCharacter",
		routeNumber: "162V",
		routeType:   apiv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "Vinbus",
		routeNumber: "D4",
		routeType:   apiv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "IntercityBus",
		routeNumber: "60-01",
		routeType:   apiv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "TourBus",
		routeNumber: "DL01",
		routeType:   apiv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "HCMCMetro",
		routeNumber: "MRT1",
		routeType:   apiv1.RouteType_ROUTE_TYPE_METRO,
	}, {
		name:        "SaigonWaterbus",
		routeNumber: "SWB1",
		routeType:   apiv1.RouteType_ROUTE_TYPE_FERRY,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			result, err := routeType(&gobus.Route{Number: test.routeNumber})
			assert.NoError(err)
			assert.Equal(test.routeType, result)
		})
	}
}
