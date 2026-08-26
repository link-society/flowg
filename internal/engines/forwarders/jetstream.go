package forwarders

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
	"link-society.com/flowg/internal/models"
)

// jetstreamRuntime sends records to nats JetStream server
type jetstreamRuntime struct {
	config *models.ForwarderJetstreamV2

	client     jetstream.JetStream
	connection *nats.Conn

	body      *vm.Program
	subject   *vm.Program
	messageID *vm.Program
}

var _ Runtime = (*jetstreamRuntime)(nil)

func (rt *jetstreamRuntime) Init(ctx context.Context) error {
	reply := make(chan error, 1)
	defer close(reply)

	go func() {
		var err error

		opts := nats.Options{
			Servers: rt.config.Servers,
			Secure:  !rt.config.TlsInsecureSkipVerify,
			TLSConfig: new(tls.Config{
				ServerName:         rt.config.TlsServerName,
				InsecureSkipVerify: rt.config.TlsInsecureSkipVerify,
			}),
			Timeout: 0,
			Token:   rt.config.Token,
		}

		switch rt.config.AuthMode {
		case "none":
		case "user_password":
			opts.User = rt.config.Username
			opts.Password = rt.config.Password
		case "token":
			opts.Token = rt.config.Token
		case "credentials":
			userJWT, err := jwt.ParseDecoratedJWT([]byte(rt.config.Credentials))
			if err != nil {
				reply <- fmt.Errorf("can't parse JWT credentials: %w", err)
			}

			keyPair, err := nkeys.ParseDecoratedNKey([]byte(rt.config.Credentials))
			if err != nil {
				reply <- fmt.Errorf("can't parse JWT key pair: %w", err)
			}

			opts.UserJWT = func() (string, error) {
				return userJWT, nil
			}

			opts.SignatureCB = func(nonce []byte) ([]byte, error) {
				return keyPair.Sign(nonce)
			}
		default:
			reply <- fmt.Errorf("unknown auth_mode: %s", rt.config.AuthMode)
		}

		if len(rt.config.TlsCertificate) > 0 {
			clientCAs := x509.NewCertPool()
			if !clientCAs.AppendCertsFromPEM([]byte(rt.config.TlsCertificate)) {
				reply <- fmt.Errorf("couldn't parse client certificate")
			}

			opts.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
			opts.TLSConfig.ClientCAs = clientCAs
		}

		if len(rt.config.TlsCA) > 0 {
			clientCert, err := tls.X509KeyPair(
				[]byte(rt.config.TlsCA),
				[]byte(rt.config.TlsPrivateKey),
			)
			if err != nil {
				reply <- fmt.Errorf("couldn't parse PEM CA certificate: %w", err)
			}
			opts.TLSConfig.Certificates = append(opts.TLSConfig.Certificates, clientCert)
		}

		rt.connection, err = opts.Connect()
		if err != nil {
			reply <- fmt.Errorf("failed to connect to nats server: %w", err)
		}

		rt.client, err = jetstream.New(rt.connection)
		if err != nil {
			reply <- fmt.Errorf("failed to establish jetstream connection: %w", err)
		}

		body := rt.config.Body
		if body == "" {
			body = "@expr:toJSON(log)"
		}

		rt.body, err = CompileDynamicField(string(body))
		if err != nil {
			reply <- fmt.Errorf("failed to compile body field: %w", err)
		}

		rt.subject, err = CompileDynamicField(string(rt.config.Subject))
		if err != nil {
			reply <- fmt.Errorf("failed to compile subject field: %w", err)
		}

		rt.messageID, err = CompileDynamicField(string(rt.config.MessageID))
		if err != nil {
			reply <- fmt.Errorf("failed to compile message_id field: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-reply:
		return err
	}
}

func (rt *jetstreamRuntime) Close(context.Context) error {
	rt.connection.Close()

	return nil
}

func (rt *jetstreamRuntime) Call(ctx context.Context, record *models.LogRecord) error {
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

	body, err := eval(rt.body, "body")
	if err != nil {
		return fmt.Errorf("failed to evaluate `body` record: %w", err)
	}

	subject, err := eval(rt.subject, "subject")
	if err != nil {
		return fmt.Errorf("failed to evaluate `subject` record: %w", err)
	}

	messageID, err := eval(rt.messageID, "message_id")
	if err != nil {
		return fmt.Errorf("failed to evaluate `message_id` record: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(rt.config.PublishTimeout)*time.Second)
	defer cancel()

	streamName, err := rt.client.StreamNameBySubject(ctx, "orders.created")
	if err != nil {
		return err
	}

	if streamName != rt.config.ExpectedStream {
		return fmt.Errorf("expected stream does not match actual stream")
	}

	_, err = rt.client.PublishMsg(ctx, &nats.Msg{
		Subject: subject,
		Data:    []byte(body),
		Header:  rt.config.Headers,
	}, jetstream.WithMsgID(messageID))
	if err != nil {
		return fmt.Errorf("failed to publish record: %w", err)
	}

	return err
}
