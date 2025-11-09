package models

import (
	"context"
	"fmt"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type forwarderStateAmqpV2 struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type ForwarderAmqpV2 struct {
	Type       string `json:"type" enum:"amqp" required:"true"`
	Url        string `json:"url" required:"true"`
	Exchange   string `json:"exchange" required:"true"`
	RoutingKey string `json:"routing_key" required:"true"`

	state *forwarderStateAmqpV2
}

func (f *ForwarderAmqpV2) init(ctx context.Context) error {
	reply := make(chan error, 1)
	defer close(reply)

	go func() {
		conn, err := amqp.Dial(f.Url)
		if err != nil {
			reply <- fmt.Errorf("failed to connect to AMQP server: %w", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			reply <- fmt.Errorf("failed to open an AMQP channel: %w", err)
		}

		f.state = &forwarderStateAmqpV2{
			conn:    conn,
			channel: ch,
		}

		reply <- nil
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-reply:
		return err
	}
}

func (f *ForwarderAmqpV2) close(context.Context) error {
	if f.state != nil {
		if f.state.channel != nil {
			f.state.channel.Close()
		}
		if f.state.conn != nil {
			f.state.conn.Close()
		}
	}
	return nil
}

func (f *ForwarderAmqpV2) call(ctx context.Context, record *LogRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	err = f.state.channel.PublishWithContext(
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
