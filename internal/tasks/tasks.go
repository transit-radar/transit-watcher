package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskScheduleMultiGoGeolocation = "task:schedule:multigo:geolocation"

	TaskAggregateGoBus              = "task:aggregate:gobus"
	TaskAggregateGoBusStops         = "task:aggregate:gobus:stop"
	TaskAggregateMultiGoGeolocation = "task:aggregate:multigo:geolocation"
)

func NewScheduleMultiGoGeolocation() (*asynq.Task, error) {
	return asynq.NewTask(TaskScheduleMultiGoGeolocation, nil), nil
}

func NewAggregateGoBusTask() (*asynq.Task, error) {
	return asynq.NewTask(TaskAggregateGoBus, nil), nil
}

func NewAggregateGoBusStopsTask() (*asynq.Task, error) {
	return asynq.NewTask(TaskAggregateGoBusStops, nil), nil
}

type AggregateMultiGoGeolocationParams struct {
	RouteID   string
	VariantID string
	Direction int
}

func NewAggregateMultiGoGeolocationTask(params AggregateMultiGoGeolocationParams) (*asynq.Task, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskAggregateMultiGoGeolocation, payload), nil
}
