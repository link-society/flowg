package models

// ForwarderAwsCloudWatchV2 forwards records to an AWS CloudWatch Logs stream,
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
}
