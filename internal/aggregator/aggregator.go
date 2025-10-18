package aggregator

import (
	"context"

	"buf.build/gen/go/catou/transit-radar/connectrpc/go/api/v1/apiv1connect"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/catouberos/transit-watcher/providers/multigo"
)

type Aggregator interface {
	Aggregate(context.Context) error
}

type aggregator struct {
	routeService       apiv1connect.RouteServiceClient
	variantService     apiv1connect.VariantServiceClient
	geolocationService apiv1connect.GeolocationServiceClient
	stopService        apiv1connect.StopServiceClient

	goBusClient   *gobus.Client
	multiGoClient *multigo.Client
}

func NewAggregator(
	routeService apiv1connect.RouteServiceClient,
	variantService apiv1connect.VariantServiceClient,
	geolocationService apiv1connect.GeolocationServiceClient,
	stopService apiv1connect.StopServiceClient,
	goBusClient *gobus.Client,
	multiGoClient *multigo.Client,
) Aggregator {
	return &aggregator{
		routeService:       routeService,
		variantService:     variantService,
		geolocationService: geolocationService,
		stopService:        stopService,

		goBusClient:   goBusClient,
		multiGoClient: multiGoClient,
	}
}
