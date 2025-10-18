package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"buf.build/gen/go/catou/transit-radar/connectrpc/go/api/v1/apiv1connect"
	"github.com/catouberos/transit-watcher/internal/aggregator"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/catouberos/transit-watcher/providers/multigo"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	routeService := apiv1connect.NewRouteServiceClient(http.DefaultClient, "http://localhost:5001")
	variantService := apiv1connect.NewVariantServiceClient(http.DefaultClient, "http://localhost:5001")
	geolocationService := apiv1connect.NewGeolocationServiceClient(http.DefaultClient, "http://localhost:5001")

	client := &http.Client{
		Timeout: 2 * time.Minute,
	}

	agg := aggregator.NewAggregator(
		routeService,
		variantService,
		geolocationService,
		gobus.NewClient(client),
		multigo.NewClient(client),
	)
	agg.Aggregate(context.Background())

	<-ctx.Done()
	slog.Info("attempt to gracefully shutdown")
}
