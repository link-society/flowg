package models

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type forwarderStateAwsCloudWatchV2 struct {
	client *cloudwatchlogs.Client
}

// ForwarderCloudWatchV2 forwards records to an AWS CloudWatch Logs stream,
// authenticating with static credentials.
type ForwarderAwsCloudWatchV2 struct {
	Type     string `json:"type" enum:"awscloudwatch" required:"true"`
	AppID    string `json:"app_id"`
	Endpoint string `json:"endpoint" required:"true"`

	Region string `json:"region"`

	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`

	Group  string `json:"group" required:"true"`
	Stream string `json:"stream" required:"true"`

	state *forwarderStateAwsCloudWatchV2
}

func (f *ForwarderAwsCloudWatchV2) init(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(f.AccessKeyID, f.SecretAccessKey, f.SessionToken),
		),
	)

	if err != nil {
		return fmt.Errorf("failed to acquire credentials: %w", err)
	}

	f.state = &forwarderStateAwsCloudWatchV2{
		client: cloudwatchlogs.New(cloudwatchlogs.Options{
			AppID:        f.AppID,
			BaseEndpoint: &f.Endpoint,
			Credentials:  cfg.Credentials,
			Region:       f.Region,
		}),
	}

	return nil
}

func (f *ForwarderAwsCloudWatchV2) close(context.Context) error {
	return nil
}

func (f *ForwarderAwsCloudWatchV2) call(ctx context.Context, record *LogRecord) error {
	message, err := json.Marshal(record.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	event := types.InputLogEvent{
		Message:   new(string(message)),
		Timestamp: new(record.Timestamp.Unix()),
	}

	_, err = f.state.client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     []types.InputLogEvent{event},
		LogGroupName:  &f.Group,
		LogStreamName: &f.Stream,
	})

	return err
}
