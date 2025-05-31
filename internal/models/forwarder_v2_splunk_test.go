package models_test

import (
	"net/http/httptest"
	"testing"

	"context"
	"net/http"

	"link-society.com/flowg/internal/models"
)

func TestForwarderSplunk_Call(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		auth := r.Header.Get("Authorization")
		if auth != "Splunk test-token" {
			t.Fatalf("unexpected auth header: %s", auth)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Fatalf("unexpected content type: %s", contentType)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"text": "Success", "code": 0}`))
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Splunk: &models.ForwarderSplunkV2{
				Endpoint: testServer.URL,
				Token:    "test-token",
			},
		},
	}

	record := models.NewLogRecord(map[string]string{
		"message": "test message",
		"host":    "test-host",
	})
	err := forwarder.Call(context.Background(), record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForwarderSplunk_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"text": "Internal Server Error", "code": 1}`))
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Splunk: &models.ForwarderSplunkV2{
				Endpoint: testServer.URL,
				Token:    "test-token",
			},
		},
	}

	record := models.NewLogRecord(map[string]string{})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error")
	}
}
