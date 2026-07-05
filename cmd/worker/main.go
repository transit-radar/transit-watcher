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
	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
	"codeberg.org/transit-radar/transit-watcher/provider/multigo"
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

	clients, err := clients.InitClients(&config)
	if err != nil {
		slog.Error("cannot load clients", "error", err)
		return err
	}
	defer clients.Close()

	a := aggregator.NewAggregator(
		clients.Kafka,
		clients.Redis,
		gobus.NewClient(),
		multigo.NewClient(),
	)

	server := asynq.NewServer(asynq.RedisClientOpt{
		Addr: config.Redis.Address,
	}, asynq.Config{
		Concurrency: 2,
	})

	mux := asynq.NewServeMux()
	mux.Handle(tasks.TaskScheduleMultiGoGeolocation, tasks.NewScheduleMultiGoProcessorHandler(clients))
	mux.Handle(tasks.TaskAggregateGoBus, tasks.NewAggregateGoBusHandler(a))
	mux.Handle(tasks.TaskAggregateMultiGoGeolocation, tasks.NewAggregateMultiGoGeolocationHandler(a))

	go func() {
		<-ctx.Done()
		server.Shutdown()
	}()

	if err := server.Run(mux); err != nil {
		return err
	}

	return nil
}
