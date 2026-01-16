package models

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	amqp "github.com/rabbitmq/amqp091-go"
)

type forwarderStateAmqpV2 struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	exchange   *vm.Program
	routingKey *vm.Program
	body       *vm.Program
}

type ForwarderAmqpV2 struct {
	Type       string                         `json:"type" enum:"amqp" required:"true"`
	Url        string                         `json:"url" required:"true" format:"uri"`
	Exchange   ForwarderAmqpV2ExchangeField   `json:"exchange" required:"true" minLength:"1"`
	RoutingKey ForwarderAmqpV2RoutingKeyField `json:"routing_key" default:""`
	Body       ForwarderAmqpV2BodyField       `json:"body,omitempty"`

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
			_ = conn.Close()
			reply <- fmt.Errorf("failed to open an AMQP channel: %w", err)
		}

		f.state = &forwarderStateAmqpV2{
			conn:    conn,
			channel: ch,
		}

		f.state.exchange, err = CompileDynamicField(string(f.Exchange))
		if err != nil {
			reply <- fmt.Errorf("failed to compile exchange field: %w", err)
		}

		f.state.routingKey, err = CompileDynamicField(string(f.RoutingKey))
		if err != nil {
			reply <- fmt.Errorf("failed to compile routingKey field: %w", err)
		}

		body := f.Body
		if body == "" {
			body = "@expr:toJSON(log)"
		}
		f.state.body, err = CompileDynamicField(string(body))
		if err != nil {
			reply <- fmt.Errorf("failed to compile body field: %w", err)
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
			_ = f.state.channel.Close()
		}
		if f.state.conn != nil {
			_ = f.state.conn.Close()
		}
	}
	return nil
}

func (f *ForwarderAmqpV2) call(ctx context.Context, record *LogRecord) error {
	env := map[string]any{
		"timestamp": record.Timestamp,
		"log":       record.Fields,
	}

	eval := func(prog *vm.Program, field string) (string, error) {
		out, err := expr.Run(prog, env)
		if err != nil {
			return "", fmt.Errorf("failed to evaluate %s expression: %w", field, err)
		}
		str, ok := out.(string)
		if !ok {
			return "", fmt.Errorf("%s expression did not evaluate to string", field)
		}
		return str, nil
	}

	exchange, err := eval(f.state.exchange, "exchange")
	if err != nil {
		return fmt.Errorf("failed to evaluate `exchange` record: %w", err)
	}

	routingKey, err := eval(f.state.routingKey, "routingKey")
	if err != nil {
		return fmt.Errorf("failed to evaluate `routingKey` record: %w", err)
	}

	body, err := eval(f.state.body, "body")
	if err != nil {
		return fmt.Errorf("failed to evaluate `body` record: %w", err)
	}

	rec := map[string]any{
		"timestamp": record.Timestamp,
		"body":      body,
	}

	payload, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	err = f.state.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
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
