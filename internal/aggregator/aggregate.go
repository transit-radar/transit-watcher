package aggregator

import (
	"context"
	"log/slog"
	"strings"

	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	"github.com/catouberos/transit-watcher/internal/utils"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"github.com/cenkalti/backoff/v5"
	"github.com/gohugoio/hashstructure"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

const cutset = " \n\r	"

func (a *aggregator) Aggregate(ctx context.Context) error {
	// get routes, variants, and variants-stops
	routes, err := backoff.Retry(
		ctx,
		func() ([]gobus.Route, error) {
			slog.InfoContext(ctx, "getting routes...")
			return a.goBusClient.GetRoutes(ctx)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(5),
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot get routes", "error", err)
		return err
	}
	slog.InfoContext(ctx, "get routes successfully", "count", len(routes))

	for _, route := range routes {
		hash, err := hashstructure.Hash(route, nil)
		if err != nil {
			slog.ErrorContext(ctx, "cannot calculate hash", "error", err)
			continue
		}

		routeType, err := utils.RouteType(&route)
		if err != nil {
			slog.ErrorContext(ctx, "cannot parse route type", "error", err)
			continue
		}

		r := &processorv1beta1.Route{
			Id: &processorv1beta1.Identity{
				Identifier: processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_EBMS,
				Value:      string(route.Id),
			},
			Number: route.Number,
			Name:   strings.Trim(route.Name, cutset),
			Type:   routeType,
			Hash:   hash,
		}

		b, err := proto.Marshal(r)
		if err != nil {
			slog.ErrorContext(ctx, "cannot marshal", "error", err)
			continue
		}

		record := &kgo.Record{Topic: "processor.route.update", Value: b}
		a.kafka.Produce(ctx, record, func(_ *kgo.Record, err error) {
			if err != nil {
				slog.ErrorContext(ctx, "record had a produce error", "error", err)
			} else {
				slog.InfoContext(ctx, "send success!")
			}
		})

		for _, variant := range route.Variants {
			hash, err := hashstructure.Hash(variant, nil)
			if err != nil {
				slog.ErrorContext(ctx, "cannot calculate hash", "error", err)
				continue
			}

			v := &processorv1beta1.Variant{
				Id: &processorv1beta1.Identity{
					Identifier: processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_EBMS,
					Value:      string(variant.Id),
				},
				RouteId: &processorv1beta1.Identity{
					Identifier: processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_EBMS,
					Value:      string(route.Id),
				},
				Name:      strings.Trim(variant.Name, cutset),
				ShortName: &variant.ShortName,
				Distance:  &variant.Distance,
				Duration:  &variant.Duration,
				Hash:      hash,
			}

			b, err := proto.Marshal(v)
			if err != nil {
				slog.ErrorContext(ctx, "cannot marshal", "error", err)
				continue
			}

			record := &kgo.Record{Topic: "processor.variant.update", Value: b}
			a.kafka.Produce(ctx, record, func(_ *kgo.Record, err error) {
				if err != nil {
					slog.ErrorContext(ctx, "record had a produce error", "error", err)
				} else {
					slog.InfoContext(ctx, "send success!")
				}
			})
		}
	}

	// get stops
	stops, err := backoff.Retry(
		ctx,
		func() ([]gobus.Stop, error) {
			slog.InfoContext(ctx, "getting stops...")
			return a.goBusClient.GetStops(ctx)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(5),
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot get stops", "error", err)
		return err
	}
	slog.InfoContext(ctx, "get stops successfully", "count", len(stops))

	for _, stop := range stops {
		hash, err := hashstructure.Hash(stop, nil)
		if err != nil {
			slog.ErrorContext(ctx, "cannot calculate hash", "error", err)
			continue
		}

		s := &processorv1beta1.Stop{
			Id: &processorv1beta1.Identity{
				Identifier: processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_EBMS,
				Value:      string(stop.Property.Id),
			},
			Code: strings.Trim(stop.Property.Code, cutset),
			Name: strings.Trim(stop.Property.Name, cutset),
			Hash: hash,
		}

		b, err := proto.Marshal(s)
		if err != nil {
			slog.ErrorContext(ctx, "cannot marshal", "error", err)
			continue
		}

		record := &kgo.Record{Topic: "processor.stop.update", Value: b}
		a.kafka.Produce(ctx, record, func(_ *kgo.Record, err error) {
			if err != nil {
				slog.ErrorContext(ctx, "record had a produce error", "error", err)
			} else {
				slog.InfoContext(ctx, "send success!")
			}
		})
	}

	return nil
}
