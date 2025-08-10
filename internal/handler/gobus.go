package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-watcher/internal/crawler"
	"github.com/catouberos/transit-watcher/internal/models"
	"github.com/wagslane/go-rabbitmq"
)

func GoBusDataHandler(conn *rabbitmq.Conn, responses <-chan *crawler.CrawlerResponse, results chan []string) error {
	routePub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("route"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	variantPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("variant"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	variantStopPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("variantstop"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	for response := range responses {
		// unmarshal
		routes := []models.GoBusRoute{}
		variants := []models.GoBusRouteVariantWithDescription{}
		variantStops := []dto.VariantStopByEbmsIDImport{}

		err := json.Unmarshal(response.Body, &routes)
		if err != nil {
			continue
		}

		for _, route := range routes {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logger := slog.New(slog.Default().Handler())
			logger.With("route", route)

			params, err := NewRouteInsertData(&route)
			if err != nil {
				logger.Error("Error parsing route data", "error", err)
			}

			data, err := json.Marshal(params)
			if err != nil {
				logger.Error("Error marshal route data", "error", err)
			}

			err = routePub.PublishWithContext(
				ctx,
				data,
				[]string{"route.event.updated"},
				rabbitmq.WithPublishOptionsContentType("application/json"),
				rabbitmq.WithPublishOptionsMandatory,
				rabbitmq.WithPublishOptionsPersistentDelivery,
				rabbitmq.WithPublishOptionsExchange("route"))
			if err != nil {
				logger.Error("Cannot publish route update", "error", err)
			}

			for _, variant := range route.Variants {
				var description string

				if variant.IsOutbound {
					description = route.Info.OutboundDescription
				} else {
					description = route.Info.InboundDescription
				}

				variants = append(variants, models.GoBusRouteVariantWithDescription{
					GoBusRouteVariant: variant,
					Description:       description,
				})
			}
		}

		<-time.After(1 * time.Second)

		for _, variant := range variants {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logger := slog.New(slog.Default().Handler())
			logger.With("variant", variant)

			params, err := NewVariantInsertData(&variant)
			if err != nil {
				logger.Error("Error parsing variant data", "error", err)
				continue
			}

			data, err := json.Marshal(params)
			if err != nil {
				logger.Error("Error marshal variant data", "error", err)
				continue
			}

			err = variantPub.PublishWithContext(
				ctx,
				data,
				[]string{"variant.event.updated"},
				rabbitmq.WithPublishOptionsContentType("application/json"),
				rabbitmq.WithPublishOptionsMandatory,
				rabbitmq.WithPublishOptionsPersistentDelivery,
				rabbitmq.WithPublishOptionsExchange("variant"))
			if err != nil {
				logger.Error("Cannot publish variant update", "error", err)
			}

			for i, stop := range variant.Stops {
				routeId, err := strconv.ParseInt(variant.RouteId, 10, 64)
				if err != nil {
					logger.Error("Error parsing variant ID", "error", err)
					continue
				}

				variantId, err := strconv.ParseInt(variant.Id, 10, 64)
				if err != nil {
					logger.Error("Error parsing variant ID", "error", err)
					continue
				}

				stopId, err := strconv.ParseInt(stop.Id, 10, 64)
				if err != nil {
					logger.Error("Error parsing stop ID", "error", err)
					continue
				}

				variantStops = append(variantStops, dto.VariantStopByEbmsIDImport{
					RouteEbmsID:   routeId,
					VariantEbmsID: variantId,
					StopEbmsID:    stopId,
					OrderScore:    int32(i),
				})
			}
		}

		urls := FilterTransitRoutes(&routes)

		slog.Info("Updated routes/variants from GoBus", "count", len(urls))

		results <- urls

		<-time.After(1 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		variantStopData, err := json.Marshal(variantStops)
		if err != nil {
			slog.Error("Error marshalling variant stop data", "error", err)
			continue
		}

		err = variantStopPub.PublishWithContext(
			ctx,
			variantStopData,
			[]string{"variantstop.action.import"},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsMandatory,
			rabbitmq.WithPublishOptionsPersistentDelivery,
			rabbitmq.WithPublishOptionsExchange("variantstop"))
		if err != nil {
			slog.Error("cannot publish variant stop import", "error", err)
		}
	}

	return nil
}

func GoBusStopHandler(conn *rabbitmq.Conn, responses <-chan *crawler.CrawlerResponse) error {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("stop"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	for response := range responses {
		// unmarshal
		stops := models.GoBusStopResponse{}

		err := json.Unmarshal(response.Body, &stops)
		if err != nil {
			continue
		}

		for _, stop := range stops.Stops {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logger := slog.New(slog.Default().Handler())
			logger.With("stop", stop)

			params, err := NewStopImportData(&stop)
			if err != nil {
				logger.Error("Error parsing stop data", "error", err)
			}

			data, err := json.Marshal(params)
			if err != nil {
				logger.Error("Error marshal stop data", "error", err)
			}

			err = publisher.PublishWithContext(
				ctx,
				data,
				[]string{"stop.action.import"},
				rabbitmq.WithPublishOptionsContentType("application/json"),
				rabbitmq.WithPublishOptionsMandatory,
				rabbitmq.WithPublishOptionsPersistentDelivery,
				rabbitmq.WithPublishOptionsExchange("stop"))
			if err != nil {
				logger.Error("Cannot publish stop import", "error", err)
			}
		}
	}

	return nil
}
