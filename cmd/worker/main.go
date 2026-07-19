package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"codeberg.org/transit-radar/transit-watcher/internal/aggregator"
	"codeberg.org/transit-radar/transit-watcher/internal/clients"
	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/tasks"
	"codeberg.org/transit-radar/transit-watcher/pkg/otel"
	"github.com/hibiken/asynq"
)

const PackageName = "codeberg.org/transit-radar/transit-watcher/cmd/worker"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "unrecoverable error", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "attempt to gracefully shutdown")
}

func run(ctx context.Context) error {
	cfg := config.NewWorkerConfig()
	if err := cfg.Load(); err != nil {
		slog.ErrorContext(ctx, "cannot load application config", "error", err)
		return err
	}

	if err := otel.Init(ctx, cfg.Application.Name, PackageName); err != nil {
		slog.ErrorContext(ctx, "cannot load otel config", "error", err)
		return err
	}
	defer otel.Shutdown(context.Background())

	clients, err := clients.InitClients(cfg)
	if err != nil {
		slog.Error("cannot load clients", "error", err)
		return err
	}
	defer clients.Close()

	a, err := aggregator.NewAggregator(cfg, clients.Kafka, clients.Redis)
	if err != nil {
		slog.ErrorContext(ctx, "cannot create aggregator", "error", err)
		return err
	}

	server := asynq.NewServer(asynq.RedisClientOpt{
		Addr: cfg.Redis.Address,
	}, asynq.Config{
		Concurrency: 20,
	})

	mux := asynq.NewServeMux()
	mux.Handle(tasks.TaskScheduleMultiGoGeolocation, tasks.NewScheduleMultiGoGeolocationHandler(clients))
	mux.Handle(tasks.TaskScheduleTTGTGeolocation, tasks.NewScheduleTTGTGeolocationHandler(clients))
	mux.Handle(tasks.TaskAggregateGoBus, tasks.NewAggregateGoBusHandler(a))
	mux.Handle(tasks.TaskAggregateGoBusStops, tasks.NewAggregateGoBusStopsHandler(a))
	mux.Handle(tasks.TaskAggregateMultiGoGeolocation, tasks.NewAggregateMultiGoGeolocationHandler(a))
	mux.Handle(tasks.TaskAggregateTTGTGeolocation, tasks.NewAggregateTTGTGeolocationHandler(a))

	go func() {
		<-ctx.Done()
		server.Shutdown()
	}()

	if err := server.Run(mux); err != nil {
		slog.ErrorContext(ctx, "cannot load server", "error", err)
		return err
	}

	return nil
}
