package forwarders

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"link-society.com/flowg/internal/models"
)

// awsCloudWatchRuntime sends records as log events to an AWS CloudWatch Logs
// stream.
type awsCloudWatchRuntime struct {
	config *models.ForwarderAwsCloudWatchV2

	client *cloudwatchlogs.Client
}

var _ Runtime = (*awsCloudWatchRuntime)(nil)

func (rt *awsCloudWatchRuntime) Init(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				rt.config.AccessKeyID,
				rt.config.SecretAccessKey,
				rt.config.SessionToken,
			),
		),
	)

	if err != nil {
		return fmt.Errorf("failed to acquire credentials: %w", err)
	}

	rt.client = cloudwatchlogs.New(cloudwatchlogs.Options{
		AppID:        rt.config.AppID,
		BaseEndpoint: &rt.config.Endpoint,
		Credentials:  cfg.Credentials,
		Region:       rt.config.Region,
	})

	return nil
}

func (rt *awsCloudWatchRuntime) Close(context.Context) error {
	return nil
}

func (rt *awsCloudWatchRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	message, err := json.Marshal(record.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	event := types.InputLogEvent{
		Message:   new(string(message)),
		Timestamp: new(record.Timestamp.Unix()),
	}

	_, err = rt.client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     []types.InputLogEvent{event},
		LogGroupName:  &rt.config.Group,
		LogStreamName: &rt.config.Stream,
	})

	return err
}
