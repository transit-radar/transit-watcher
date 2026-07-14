package aggregator

import (
	"time"

	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/events"
	"codeberg.org/transit-radar/transit-watcher/internal/processor"
	"codeberg.org/transit-radar/transit-watcher/internal/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/store"
	"codeberg.org/transit-radar/transit-watcher/provider/ebms"
	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
	"codeberg.org/transit-radar/transit-watcher/provider/multigo"
	"codeberg.org/transit-radar/transit-watcher/provider/ttgt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type Aggregator struct {
	route       processor.RouteProcessor
	variant     processor.VariantProcessor
	geolocation processor.GeolocationProcessor

	ebms    ebms.Client
	goBus   gobus.Client
	multiGo multigo.Client
	ttgt    ttgt.Client

	rateLimit *rate.Limiter
}

func NewAggregator(
	config *config.Config,
	kafka sarama.Client,
	redis *redis.Client,
) *Aggregator {
	eventHandler, err := events.NewKafkaEventHandler(config.Application.Name, kafka)
	if err != nil {
		panic(err)
	}

	store := store.NewRedisStore(redis)

	route := v1beta1.NewRouteProcessor(config, eventHandler, store)
	variant := v1beta1.NewVariantProcessor(config, eventHandler, store)
	geolocation := v1beta1.NewGeolocationProcessor(config, eventHandler, store)

	return &Aggregator{
		// processors
		route:       route,
		variant:     variant,
		geolocation: geolocation,

		// clients
		goBus:   gobus.NewClient(),
		multiGo: multigo.NewClient(),
		ttgt:    ttgt.NewClient(),

		// currently limits 1 event per 50ms
		rateLimit: rate.NewLimiter(rate.Every(50*time.Millisecond), 1),
	}
}
