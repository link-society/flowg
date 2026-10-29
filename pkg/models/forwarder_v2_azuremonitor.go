package models

// ForwarderAzureMonitorV2 forwards records to the Microsoft Azure Monitor Logs stream,
// authenticating with a static token with expiry time.
type ForwarderAzureMonitorV2 struct {
	Type          string `json:"type" enum:"azuremonitor" required:"true"`
	Endpoint      string `json:"endpoint"`
	Token         string `json:"token"`
	ExpiresOn     string `json:"expires_on"`
	RuleID        string `json:"rule_id"`
	StreamName    string `json:"stream_name"`
	AllowInsecure bool   `json:"allow_insecure" default:"false"`
}
