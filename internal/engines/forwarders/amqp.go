package forwarders

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	amqp "github.com/rabbitmq/amqp091-go"

	"link-society.com/flowg/internal/models"
)

type amqpRuntime struct {
	config *models.ForwarderAmqpV2

	conn    *amqp.Connection
	channel *amqp.Channel

	exchange   *vm.Program
	routingKey *vm.Program
	body       *vm.Program
}

var _ Runtime = (*amqpRuntime)(nil)

func (rt *amqpRuntime) Init(ctx context.Context) error {
	reply := make(chan error, 1)
	defer close(reply)

	go func() {
		conn, err := amqp.Dial(rt.config.Url)
		if err != nil {
			reply <- fmt.Errorf("failed to connect to AMQP server: %w", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			_ = conn.Close()
			reply <- fmt.Errorf("failed to open an AMQP channel: %w", err)
		}

		rt.conn = conn
		rt.channel = ch

		rt.exchange, err = CompileDynamicField(string(rt.config.Exchange))
		if err != nil {
			reply <- fmt.Errorf("failed to compile exchange field: %w", err)
		}

		rt.routingKey, err = CompileDynamicField(string(rt.config.RoutingKey))
		if err != nil {
			reply <- fmt.Errorf("failed to compile routingKey field: %w", err)
		}

		body := rt.config.Body
		if body == "" {
			body = "@expr:toJSON(log)"
		}
		rt.body, err = CompileDynamicField(string(body))
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

func (rt *amqpRuntime) Close(context.Context) error {
	if rt != nil {
		if rt.channel != nil {
			_ = rt.channel.Close()
		}
		if rt.conn != nil {
			_ = rt.conn.Close()
		}
	}
	return nil
}

func (rt *amqpRuntime) Call(ctx context.Context, record *models.LogRecord) error {
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

	exchange, err := eval(rt.exchange, "exchange")
	if err != nil {
		return fmt.Errorf("failed to evaluate `exchange` record: %w", err)
	}

	routingKey, err := eval(rt.routingKey, "routingKey")
	if err != nil {
		return fmt.Errorf("failed to evaluate `routingKey` record: %w", err)
	}

	body, err := eval(rt.body, "body")
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

	err = rt.channel.PublishWithContext(
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
