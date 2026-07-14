package events

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
)

type kafkaEventHandler struct {
	name string

	kafka sarama.Client
}

func NewKafkaEventHandler(name string, kafka sarama.Client) (*kafkaEventHandler, error) {
	return &kafkaEventHandler{name: name, kafka: kafka}, nil
}

func (e *kafkaEventHandler) Send(ctx context.Context, topic string, event cloudevents.Event, opts ...client.Option) error {
	sender, err := kafka_sarama.NewSenderFromClient(e.kafka, topic)
	if err != nil {
		return err
	}
	defer sender.Close(ctx)

	c, err := cloudevents.NewClient(sender, opts...)
	if err != nil {
		return err
	}

	return c.Send(ctx, event)
}
