package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/scheduler"
	"codeberg.org/transit-radar/transit-watcher/pkg/otel"
)

const PackageName = "codeberg.org/transit-radar/transit-watcher/cmd/scheduler"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "unrecoverable error occurred", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "attempt to gracefully shutdown")
}

func run(ctx context.Context) error {
	cfg := config.NewSchedulerConfig()
	if err := cfg.Load(); err != nil {
		return fmt.Errorf("cannot load scheduler config: %w", err)
	}

	if err := otel.Init(ctx, cfg.Application.Name, PackageName); err != nil {
		return err
	}
	defer otel.Shutdown(context.Background())

	scheduler, err := scheduler.NewScheduler(ctx, cfg)
	if err != nil {
		return err
	}

	if err := scheduler.RegisterConfigured(ctx); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		scheduler.Shutdown()
	}()

	if err := scheduler.Run(); err != nil {
		return err
	}

	return nil
}
