package aggregator

import (
	"context"

	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/catouberos/transit-watcher/providers/multigo"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Aggregator interface {
	Aggregate(context.Context) error
}

type aggregator struct {
	kafka         *kgo.Client
	goBusClient   *gobus.Client
	multiGoClient *multigo.Client
}

func NewAggregator(
	kafka *kgo.Client,
	goBusClient *gobus.Client,
	multiGoClient *multigo.Client,
) Aggregator {
	return &aggregator{
		kafka:         kafka,
		goBusClient:   goBusClient,
		multiGoClient: multiGoClient,
	}
}
