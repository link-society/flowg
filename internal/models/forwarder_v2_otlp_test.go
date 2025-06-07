package models_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"link-society.com/flowg/internal/models"
)

func TestForwarderOtlpV2_Call_Success(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check content type header
		if got := r.Header.Get("Content-Type"); got != "application/x-protobuf" {
			t.Fatalf("unexpected Content-Type: %s", got)
		}

		// Check custom header
		if got := r.Header.Get("X-Test-Header"); got != "test-value" {
			t.Fatalf("unexpected X-Test-Header: %s", got)
		}

		// Read body (protobuf bytes)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		if len(body) == 0 {
			t.Fatalf("empty body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Otlp: &models.ForwarderOtlpV2{
				Type: "otlp",
				Config: models.OtlpForwarderConfig{
					Endpoint: testServer.URL,
					Headers: map[string]string{
						"X-Test-Header": "test-value",
					},
				},
			},
		},
	}

	record := models.NewLogRecord(map[string]string{
		"body": "test log body",
		"foo":  "bar",
	})
	err := forwarder.Call(context.Background(), record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForwarderOtlpV2_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return 500 error
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Otlp: &models.ForwarderOtlpV2{
				Type: "otlp",
				Config: models.OtlpForwarderConfig{
					Endpoint: testServer.URL,
				},
			},
		},
	}

	record := models.NewLogRecord(map[string]string{
		"body": "test log body",
	})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error on failure response, got nil")
	}
}
