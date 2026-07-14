package forwarders

import (
	"context"
	"fmt"

	"time"

	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"

	"link-society.com/flowg/internal/models"
)

// azureMonitorRuntime uploads records to Azure Monitor Logs, targeting the
// configured data collection rule and stream.
type azureMonitorRuntime struct {
	config *models.ForwarderAzureMonitorV2
	client *azlogs.Client
}

// staticAzureTokenCredential implements [azcore.TokenCredential] with the
// pre-issued token from the configuration instead of an interactive credential
// flow.
type staticAzureTokenCredential struct {
	token     string
	expiresOn time.Time
}

var _ Runtime = (*azureMonitorRuntime)(nil)
var _ azcore.TokenCredential = (*staticAzureTokenCredential)(nil)

// GetToken implements [azcore.TokenCredential]. It always returns the static
// token, regardless of the requested options.
func (c staticAzureTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     c.token,
		ExpiresOn: c.expiresOn,
	}, nil
}

func (rt *azureMonitorRuntime) Init(ctx context.Context) error {
	var err error
	var expires time.Time

	expires, err = time.Parse(time.RFC3339, rt.config.ExpiresOn)
	if err != nil {
		return fmt.Errorf("error parsing expiry date for Azure Monitor: %w", err)
	}

	credentials := staticAzureTokenCredential{
		token:     rt.config.Token,
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
						InsecureSkipVerify: rt.config.AllowInsecure,
					},
				},
			},
		},
	}

	rt.client, err = azlogs.NewClient(
		rt.config.Endpoint,
		credentials,
		&clientOptions,
	)
	if err != nil {
		return fmt.Errorf("failed to authenticate to Azure Monitor: %w", err)
	}

	return nil
}

func (rt *azureMonitorRuntime) Close(context.Context) error {
	return nil
}

func (rt *azureMonitorRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	records := []map[string]string{
		record.Fields,
	}

	message, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	_, err = rt.client.Upload(ctx, rt.config.RuleID, rt.config.StreamName, message, nil)
	if err != nil {
		return fmt.Errorf("failed to upload logs to Azure Monitor: %w", err)
	}

	return nil
}
