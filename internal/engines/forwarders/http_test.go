package forwarders_test

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"link-society.com/flowg/internal/engines/forwarders"
	"link-society.com/flowg/internal/models"
)

func TestForwarderHttp_Call(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fooHeader := r.Header.Get("Foo")
		if fooHeader != "Bar" {
			t.Fatalf("unexpected header value: %s", fooHeader)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Http: &models.ForwarderHttpV2{
				Url: testServer.URL,
				Headers: map[string]string{
					"Foo": "Bar",
				},
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

func TestForwarderHttp_Call_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Http: &models.ForwarderHttpV2{
				Url:     testServer.URL,
				Headers: map[string]string{},
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
