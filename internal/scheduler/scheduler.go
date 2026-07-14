package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"time"

	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Scheduler interface {
	Register(ctx context.Context, spec config.TaskSpec, task *asynq.Task) error
	RegisterConfigured(ctx context.Context) error

	Run() error
	Shutdown()
}

type scheduler struct {
	cfg            *config.SchedulerConfig
	asynqScheduler *asynq.Scheduler
}

func NewScheduler(ctx context.Context, cfg *config.SchedulerConfig) (Scheduler, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, fmt.Errorf("cannot load local location: %w", err)
	}

	asynqScheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: cfg.Redis.Address},
		&asynq.SchedulerOpts{
			Location: loc,
			PreEnqueueFunc: func(task *asynq.Task, opts []asynq.Option) {
				// otel trace propagation
				tracer := otel.GetTracerProvider().Tracer(cfg.Application.Name)
				ctx, _ := tracer.Start(ctx, fmt.Sprintf("scheduler/%s", task.Type()))

				slog.DebugContext(ctx, "scheduling task...", slog.String("task", task.Type()))

				prop := otel.GetTextMapPropagator()
				carrier := propagation.MapCarrier{}
				prop.Inject(ctx, &carrier)

				maps.Copy(task.Headers(), carrier)
			},
		},
	)

	return &scheduler{cfg, asynqScheduler}, nil
}

func (s *scheduler) Shutdown() {
	s.asynqScheduler.Shutdown()
}

func (s *scheduler) Run() error {
	return s.asynqScheduler.Run()
}
