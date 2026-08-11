---
sidebar_position: 10
---

# AWS CloudWatch

This forwarder is used to send a log record to an
[AWS CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/WhatIsCloudWatchLogs.html)
stream.

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str app_id "Identifier appended to the User-Agent of the AWS SDK"
    str endpoint "URL of the CloudWatch Logs endpoint"
    str region "AWS region to send the records to"
    str access_key_id "Access Key ID used to authenticate"
    str secret_access_key[1] "Secret Access Key used to authenticate"
    str session_token[1] "Session token, when using temporary credentials"
    str group "Name of the log group to store the log in"
    str stream "Name of the log stream to store the log in"
  }
```

:::note

1. The secret access key and the session token are **NOT** encrypted in the
   database.
2. The log group and the log stream must exist prior to sending records, they
   are not created by the forwarder.

:::

## Behavior

```go
client := cloudwatchlogs.New(cloudwatchlogs.Options{
  AppID:        appID,
  BaseEndpoint: &endpoint,
  Credentials:  credentials.NewStaticCredentialsProvider(
    accessKeyID,
    secretAccessKey,
    sessionToken,
  ),
  Region: region,
})

message, err := json.Marshal(logRecord.Fields)
// ...

_, err = client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
  LogEvents: []types.InputLogEvent{
    {
      Message:   new(string(message)),
      Timestamp: new(logRecord.Timestamp.UnixMilli()),
    },
  },
  LogGroupName:  &group,
  LogStreamName: &stream,
})
// ...
```
