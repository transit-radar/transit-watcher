package events

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func CreateEvent[T proto.Message](message T) (cloudevents.Event, error) {
	payload, err := proto.Marshal(message)
	if err != nil {
		return cloudevents.NewEvent(), err
	}

	event := cloudevents.NewEvent()
	event.SetID(uuid.NewString())
	event.SetSpecVersion("1.0")
	event.SetSource("example/uri")
	event.SetType(string(message.ProtoReflect().Descriptor().Name()))
	event.SetData("application/protobuf", payload)

	return event, nil
}
