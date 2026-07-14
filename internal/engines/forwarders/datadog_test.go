package forwarders_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"link-society.com/flowg/internal/engines/forwarders"
	"link-society.com/flowg/internal/models"
)

func TestForwarderDatadog_Call(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Datadog: &models.ForwarderDatadogV2{
				Url:    testServer.URL,
				ApiKey: "apiKey",
			},
		},
	}

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{})
	if err := runtime.Call(t.Context(), record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := runtime.Close(t.Context()); err != nil {
		t.Fatalf("failed to close forwarder: %v", err)
	}
}

func TestForwarderDatadog_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Datadog: &models.ForwarderDatadogV2{
				Url:    testServer.URL,
				ApiKey: "apiKey",
			},
		},
	}

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{})
	if err := runtime.Call(t.Context(), record); err == nil {
		t.Fatalf("expected error")
	}

	if err := runtime.Close(t.Context()); err != nil {
		t.Fatalf("failed to close forwarder: %v", err)
	}
}
