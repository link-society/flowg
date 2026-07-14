package models

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
)

type forwarderStateAzureMonitorV2 struct {
	client *azlogs.Client
}

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

	state *forwarderStateAzureMonitorV2
}

type staticTokenCredential struct {
	token     string
	expiresOn time.Time
}

func (c staticTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     c.token,
		ExpiresOn: c.expiresOn,
	}, nil
}

func (f *ForwarderAzureMonitorV2) init(ctx context.Context) error {
	var err error
	var expires time.Time

	expires, err = time.Parse(time.RFC3339, f.ExpiresOn)
	if err != nil {
		return fmt.Errorf("error parsing expiry date for Azure Monitor: %w", err)
	}

	credentials := staticTokenCredential{
		token:     f.Token,
		expiresOn: expires,
	}

	clientOptions := azlogs.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
			Transport: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: f.AllowInsecure,
					},
				},
			},
		},
	}

	f.state = &forwarderStateAzureMonitorV2{}
	f.state.client, err = azlogs.NewClient(f.Endpoint, credentials, &clientOptions)
	if err != nil {
		return fmt.Errorf("failed to authenticate to Azure Monitor: %w", err)
	}

	return nil
}

func (f *ForwarderAzureMonitorV2) close(context.Context) error {
	return nil
}

func (f *ForwarderAzureMonitorV2) call(ctx context.Context, record *LogRecord) error {
	records := []map[string]string{
		record.Fields,
	}

	message, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	_, err = f.state.client.Upload(ctx, f.RuleID, f.StreamName, message, nil)
	if err != nil {
		return fmt.Errorf("failed to upload logs to Azure Monitor: %w", err)
	}

	return nil
}
