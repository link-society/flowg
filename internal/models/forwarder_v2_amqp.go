package models

import (
	"context"
	"fmt"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ForwarderAmqpV2 struct {
	Type       string `json:"type" enum:"amqp"`
	Url        string `json:"url"`
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}

func (f *ForwarderAmqpV2) call(ctx context.Context, record *LogRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	conn, err := amqp.Dial(f.Url)
	if err != nil {
		return fmt.Errorf("failed to connect to AMQP server: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open an AMQP channel: %w", err)
	}
	defer ch.Close()

	err = ch.PublishWithContext(
		ctx,
		f.Exchange,
		f.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish AMQP message: %w", err)
	}

	return nil
}
