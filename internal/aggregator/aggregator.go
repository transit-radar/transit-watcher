package aggregator

import (
	"context"
	"sync"

	"buf.build/gen/go/catou/transit-radar/connectrpc/go/api/v1/apiv1connect"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/catouberos/transit-watcher/providers/multigo"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Aggregator interface {
	Aggregate(context.Context) error
}

type aggregator struct {
	kafka *kgo.Client

	routeService       apiv1connect.RouteServiceClient
	variantService     apiv1connect.VariantServiceClient
	geolocationService apiv1connect.GeolocationServiceClient
	stopService        apiv1connect.StopServiceClient

	goBusClient   *gobus.Client
	multiGoClient *multigo.Client

	// internal uses
	routes sync.Map
}

func NewAggregator(
	kafka *kgo.Client,
	routeService apiv1connect.RouteServiceClient,
	variantService apiv1connect.VariantServiceClient,
	geolocationService apiv1connect.GeolocationServiceClient,
	stopService apiv1connect.StopServiceClient,
	goBusClient *gobus.Client,
	multiGoClient *multigo.Client,
) Aggregator {
	return &aggregator{
		kafka: kafka,

		routeService:       routeService,
		variantService:     variantService,
		geolocationService: geolocationService,
		stopService:        stopService,

		goBusClient:   goBusClient,
		multiGoClient: multiGoClient,
	}
}
