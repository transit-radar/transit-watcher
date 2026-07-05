package aggregator

import (
	"github.com/IBM/sarama"
	"codeberg.org/transit-radar/transit-watcher/internal/events"
	"codeberg.org/transit-radar/transit-watcher/internal/processor"
	"codeberg.org/transit-radar/transit-watcher/internal/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/store"
	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
	"codeberg.org/transit-radar/transit-watcher/provider/multigo"
	"github.com/redis/go-redis/v9"
)

type Aggregator struct {
	route       processor.RouteProcessor
	variant     processor.VariantProcessor
	geolocation processor.GeolocationProcessor

	goBus   gobus.Client
	multiGo multigo.Client
}

func NewAggregator(
	kafka sarama.Client,
	redis *redis.Client,
	goBus gobus.Client,
	multiGo multigo.Client,
) *Aggregator {
	eventHandler, err := events.NewKafkaEventHandler(kafka)
	if err != nil {
		panic(err)
	}

	store := store.NewRedisStore(redis)

	route := v1beta1.NewRouteProcessor(nil, eventHandler, store)
	variant := v1beta1.NewVariantProcessor(nil, eventHandler, store)
	geolocation := v1beta1.NewGeolocationProcessor(nil, eventHandler, store)

	return &Aggregator{
		// processors
		route:       route,
		variant:     variant,
		geolocation: geolocation,

		// clients
		goBus:   goBus,
		multiGo: multiGo,
	}
}
