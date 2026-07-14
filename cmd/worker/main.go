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
	"codeberg.org/transit-radar/transit-watcher/pkg/otelhelper"
	"github.com/hibiken/asynq"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "unrecoverable error", "error", err)
		os.Exit(1)
	}

	slog.Info("attempt to gracefully shutdown")
}

func run(ctx context.Context) error {
	config, err := config.LoadConfig()
	if err != nil {
		slog.Error("cannot load application config", "error", err)
		return err
	}

	if err := otelhelper.Init(ctx, config.Application.Name); err != nil {
		return err
	}

	clients, err := clients.InitClients(&config)
	if err != nil {
		slog.Error("cannot load clients", "error", err)
		return err
	}
	defer clients.Close()

	a := aggregator.NewAggregator(
		&config,
		clients.Kafka,
		clients.Redis,
	)

	server := asynq.NewServer(asynq.RedisClientOpt{
		Addr: config.Redis.Address,
	}, asynq.Config{
		Concurrency: 20,
	})

	mux := asynq.NewServeMux()
	mux.Handle(tasks.TaskScheduleMultiGoGeolocation, tasks.NewScheduleMultiGoGeolocationHandler(clients))
	mux.Handle(tasks.TaskScheduleTTGTGeolocation, tasks.NewScheduleTTGTGeolocationHandler(clients))
	mux.Handle(tasks.TaskAggregateGoBus, tasks.NewAggregateGoBusHandler(a))
	mux.Handle(tasks.TaskAggregateMultiGoGeolocation, tasks.NewAggregateMultiGoGeolocationHandler(a))
	mux.Handle(tasks.TaskAggregateTTGTGeolocation, tasks.NewAggregateTTGTGeolocationHandler(a))

	go func() {
		<-ctx.Done()
		server.Shutdown()
	}()

	if err := server.Run(mux); err != nil {
		return err
	}

	return nil
}
