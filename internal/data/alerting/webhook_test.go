package alerting_test

import (
	"net/http/httptest"
	"testing"

	"context"
	"net/http"

	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/internal/data/alerting"
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

	webhook := &alerting.Webhook{
		Url: testServer.URL,
		Headers: map[string]string{
			"Foo": "Bar",
		},
	}

	entry := logstorage.NewLogEntry(map[string]string{})
	err := webhook.Call(context.Background(), entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhook_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	webhook := &alerting.Webhook{
		Url:     testServer.URL,
		Headers: map[string]string{},
	}

	entry := logstorage.NewLogEntry(map[string]string{})
	err := webhook.Call(context.Background(), entry)
	if err == nil {
		t.Fatalf("expected error")
	}
}
