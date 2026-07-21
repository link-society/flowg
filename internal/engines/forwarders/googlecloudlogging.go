package forwarders

import (
	"context"
	"fmt"

	"encoding/json"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"link-society.com/flowg/internal/models"
)

// googleCloudLoggingRuntime writes records to a Google Cloud Logging log,
// authenticating with a service-account JSON when one is configured.
type googleCloudLoggingRuntime struct {
	config *models.ForwarderGoogleCloudLoggingV2

	client *logging.Client
	logger *logging.Logger
}

var _ Runtime = (*googleCloudLoggingRuntime)(nil)

func (rt *googleCloudLoggingRuntime) Init(ctx context.Context) error {
	var err error

	var opts []option.ClientOption
	opts = append(opts, option.WithEndpoint(rt.config.Endpoint+":"+rt.config.EndpointPort))

	if len(rt.config.AuthJSON) > 0 {
		opts = append(opts, option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(rt.config.AuthJSON)))
	}

	if rt.config.DisableAuth {
		opts = append(opts, option.WithoutAuthentication())
	}

	if rt.config.DisableTLS {
		opts = append(opts,
			option.WithGRPCDialOption(
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			),
		)
	}

	rt.client, err = logging.NewClient(
		ctx,
		rt.config.ProjectID,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create Google Cloud client: %w", err)
	}

	rt.logger = rt.client.Logger(rt.config.LogID)

	return nil
}

func (rt *googleCloudLoggingRuntime) Close(context.Context) error {
	return rt.client.Close()
}

func (rt *googleCloudLoggingRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	message, err := json.Marshal(record.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	return rt.logger.LogSync(ctx, logging.Entry{
		Timestamp: record.Timestamp,
		Payload:   json.RawMessage(message),
	})
}
