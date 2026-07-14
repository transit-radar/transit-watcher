package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskScheduleMultiGoGeolocation = "task:schedule:multigo:geolocation"
	TaskScheduleTTGTGeolocation    = "task:schedule:ttgt:geolocation"

	TaskAggregateGoBus              = "task:aggregate:gobus"
	TaskAggregateGoBusStops         = "task:aggregate:gobus:stop"
	TaskAggregateMultiGoGeolocation = "task:aggregate:multigo:geolocation"
	TaskAggregateTTGTGeolocation    = "task:aggregate:ttgt:geolocation"
)

func NewScheduleMultiGoGeolocation() (*asynq.Task, error) {
	return asynq.NewTaskWithHeaders(TaskScheduleMultiGoGeolocation, nil, map[string]string{}), nil
}

func NewScheduleTTGTGeolocation() (*asynq.Task, error) {
	return asynq.NewTaskWithHeaders(TaskScheduleTTGTGeolocation, nil, map[string]string{}), nil
}

func NewAggregateGoBusTask() (*asynq.Task, error) {
	return asynq.NewTaskWithHeaders(TaskAggregateGoBus, nil, map[string]string{}), nil
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

	return asynq.NewTaskWithHeaders(TaskAggregateMultiGoGeolocation, payload, map[string]string{}), nil
}

type AggregateTTGTGeolocationParams struct {
	RouteID   string
	VariantID string
}

func NewAggregateTTGTGeolocationTask(params AggregateTTGTGeolocationParams) (*asynq.Task, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return asynq.NewTaskWithHeaders(TaskAggregateTTGTGeolocation, payload, map[string]string{}), nil
}
