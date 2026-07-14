package events

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"google.golang.org/protobuf/proto"
)

type EventHandler interface {
	CreateEvent(message proto.Message) (cloudevents.Event, error)
	Send(ctx context.Context, topic string, event cloudevents.Event, opts ...client.Option) error
}
