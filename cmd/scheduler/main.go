package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"codeberg.org/transit-radar/transit-watcher/internal/clients"
	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	ctx := context.Background()

	config, err := config.LoadConfig()
	if err != nil {
		slog.Error("cannot load application config", "error", err)
		os.Exit(1)
	}

	clients, err := clients.InitClients(&config)
	if err != nil {
		slog.Error("cannot load clients", "error", err)
		os.Exit(1)
	}
	defer clients.Close()

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Redis.Address})
	defer client.Close()

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		panic(err)
	}
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: config.Redis.Address},
		&asynq.SchedulerOpts{
			Location: loc,
		},
	)
	defer scheduler.Shutdown()

	if config.Task.Data.Enable {
		slog.InfoContext(ctx, "registring data task", "cron", config.Task.Data.Crontab)
		task, err := tasks.NewAggregateGoBusTask()
		if err != nil {
			slog.ErrorContext(ctx, "failed to create data task", "error", err)
			panic(err)
		}

		_, err = scheduler.Register(config.Task.Data.Crontab, task, asynq.MaxRetry(0))
		if err != nil {
			slog.ErrorContext(ctx, "failed to schedule data task", "error", err)
			panic(err)
		}
	}

	if config.Task.Geolocation.Enable {
		slog.InfoContext(ctx, "registring geolocation task", "cron", config.Task.Geolocation.Crontab)
		task, err := tasks.NewScheduleMultiGoGeolocation()
		if err != nil {
			slog.ErrorContext(ctx, "failed to create geolocation task", "error", err)
			panic(err)
		}

		_, err = scheduler.Register(config.Task.Geolocation.Crontab, task, asynq.MaxRetry(0))
		if err != nil {
			slog.ErrorContext(ctx, "failed to schedule geolocation task", "error", err)
			panic(err)
		}
	}

	if err := scheduler.Run(); err != nil {
		panic(err)
	}
}
