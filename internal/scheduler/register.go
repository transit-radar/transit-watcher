package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/tasks"
	"github.com/hibiken/asynq"
)

func (s *scheduler) RegisterConfigured(ctx context.Context) error {
	if aggregate := s.cfg.Task.Aggregate; aggregate.Enable {
		task, err := tasks.NewAggregateGoBusTask()
		if err != nil {
			return fmt.Errorf("failed to create aggregate task: %w", err)
		}

		if err := s.Register(ctx, aggregate, task); err != nil {
			return err
		}
	}

	if multigo := s.cfg.Task.MultiGoGeolocation; multigo.Enable {
		task, err := tasks.NewScheduleMultiGoGeolocation()
		if err != nil {
			return fmt.Errorf("failed to create geolocation task: %w", err)
		}

		if err := s.Register(ctx, multigo, task); err != nil {
			return err
		}
	}

	if ttgt := s.cfg.Task.TTGTGeolocation; ttgt.Enable {
		task, err := tasks.NewScheduleTTGTGeolocation()
		if err != nil {
			return fmt.Errorf("failed to create geolocation task: %w", err)
		}

		if err := s.Register(ctx, ttgt, task); err != nil {
			return err
		}
	}

	return nil
}

func (s *scheduler) Register(ctx context.Context, spec config.TaskSpec, task *asynq.Task) error {
	logger := slog.With(
		slog.String("type", task.Type()),
		slog.String("cron", spec.Crontab),
	)

	logger.InfoContext(ctx, "registring task")

	opts := []asynq.Option{}

	if spec.MaxRetry != nil {
		opts = append(opts, asynq.MaxRetry(*spec.MaxRetry))
	}

	if spec.Unique != nil {
		opts = append(opts, asynq.Unique(*spec.Unique))
	}

	id, err := s.asynqScheduler.Register(spec.Crontab, task, opts...)
	if err != nil {
		return fmt.Errorf("failed to register task: %w", err)
	}

	slog.InfoContext(ctx, "successfully registered task", slog.String("id", id))

	return nil
}
