package cmd

import "os"

var (
	defaultApiUrl     = getEnvString("FLOWG_API", "http://localhost:5080")
	defaultApiToken   = getEnvString("FLOWG_API_TOKEN", "")
	defaultMgmtApiUrl = getEnvString("FLOWG_MGMT_API", "http://localhost:9113")
)

func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
