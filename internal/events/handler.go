package events

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
)

type EventHandler interface {
	Send(ctx context.Context, topic string, event cloudevents.Event, opts ...client.Option) error
}
