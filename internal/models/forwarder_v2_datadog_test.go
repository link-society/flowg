package models_test

import (
	"net/http/httptest"
	"testing"

	"context"
	"net/http"

	"link-society.com/flowg/internal/models"
)

func TestForwarderDatadog_Call(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Datadog: &models.ForwarderDatadogV2{
				Url:    testServer.URL,
				ApiKey: "apiKey",
			},
		},
	}

	record := models.NewLogRecord(map[string]string{})
	err := forwarder.Call(context.Background(), record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForwarderDatadog_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Datadog: &models.ForwarderDatadogV2{
				Url:    testServer.URL,
				ApiKey: "apiKey",
			},
		},
	}

	record := models.NewLogRecord(map[string]string{})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error")
	}
}
