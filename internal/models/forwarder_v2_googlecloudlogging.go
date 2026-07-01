package models

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ForwarderStateGoogleCloudLoggingV2 struct {
	client *logging.Client
	logger *logging.Logger
}

type ForwarderGoogleCloudLoggingV2 struct {
	Type      string `json:"type" enum:"googlecloudlogging" required:"true"`
	Endpoint  string `json:"endpoint" required:"true"`
	ProjectID string `json:"project_id" required:"true"`
	LogID     string `json:"log_id" required:"true"`
	AuthJSON  string `json:"auth_json"`

	state *ForwarderStateGoogleCloudLoggingV2
}

func (f *ForwarderGoogleCloudLoggingV2) init(ctx context.Context) error {
	var err error

	var opts []option.ClientOption
	opts = append(opts, option.WithEndpoint(f.Endpoint))
	if len(f.AuthJSON) == 0 {
		opts = append(opts, option.WithoutAuthentication())
		opts = append(opts,
			option.WithGRPCDialOption(
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			),
		)
	} else {
		opts = append(opts, option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(f.AuthJSON)))
	}

	f.state = &ForwarderStateGoogleCloudLoggingV2{}
	f.state.client, err = logging.NewClient(
		ctx,
		f.ProjectID,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create Google Cloud client: %w", err)
	}

	f.state.logger = f.state.client.Logger(f.LogID)

	return nil
}

func (f *ForwarderGoogleCloudLoggingV2) close(context.Context) error {
	return f.state.client.Close()
}

func (f *ForwarderGoogleCloudLoggingV2) call(ctx context.Context, record *LogRecord) error {
	message, err := json.Marshal(record.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	return f.state.logger.LogSync(ctx, logging.Entry{
		Timestamp: record.Timestamp,
		Payload:   string(message),
	})
}
