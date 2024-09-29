package models_test

import (
	"net/http/httptest"
	"testing"

	"context"
	"net/http"

	"link-society.com/flowg/internal/models"
)

func TestWebhook_Call(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fooHeader := r.Header.Get("Foo")
		if fooHeader != "Bar" {
			t.Fatalf("unexpected header value: %s", fooHeader)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	webhook := &models.WebhookV1{
		Url: testServer.URL,
		Headers: map[string]string{
			"Foo": "Bar",
		},
	}

	record := models.NewLogRecord(map[string]string{})
	err := webhook.Call(context.Background(), record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhook_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	webhook := &models.WebhookV1{
		Url:     testServer.URL,
		Headers: map[string]string{},
	}

	record := models.NewLogRecord(map[string]string{})
	err := webhook.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error")
	}
}
