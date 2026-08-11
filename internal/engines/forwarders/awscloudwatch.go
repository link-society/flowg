package forwarders

import (
	"context"
	"fmt"

	"encoding/json"

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
	rt.client = cloudwatchlogs.New(cloudwatchlogs.Options{
		AppID:        rt.config.AppID,
		BaseEndpoint: &rt.config.Endpoint,
		Credentials: credentials.NewStaticCredentialsProvider(
			rt.config.AccessKeyID,
			rt.config.SecretAccessKey,
			rt.config.SessionToken,
		),
		Region: rt.config.Region,
	})

	return nil
}

func (rt *awsCloudWatchRuntime) Close(context.Context) error {
	return nil
}

func awsCloudWatchTimestamp(record *models.LogRecord) int64 {
	return record.Timestamp.UnixMilli()
}

func (rt *awsCloudWatchRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	message, err := json.Marshal(record.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	event := types.InputLogEvent{
		Message:   new(string(message)),
		Timestamp: new(awsCloudWatchTimestamp(record)),
	}

	_, err = rt.client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     []types.InputLogEvent{event},
		LogGroupName:  &rt.config.Group,
		LogStreamName: &rt.config.Stream,
	})

	return err
}
