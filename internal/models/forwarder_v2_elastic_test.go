package models_test

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"link-society.com/flowg/internal/models"
)

func TestForwarderElastic_Call_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")

		switch {
		case r.Method == "HEAD" && r.URL.Path == "/test-index":
			w.WriteHeader(http.StatusOK)

		case r.Method == "POST" && r.URL.Path == "/test-index/_doc":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"result":"created"}`))

		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))
	defer server.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	if err := forwarder.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "hello"})
	if err := forwarder.Call(t.Context(), record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := forwarder.Close(t.Context()); err != nil {
		t.Fatalf("failed to close forwarder: %v", err)
	}
}

func TestForwarderElastic_Call_IndexNotExists_CreatesIndex(t *testing.T) {
	indexCreated := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")

		switch {
		case r.Method == "HEAD" && r.URL.Path == "/test-index":
			w.WriteHeader(http.StatusNotFound)

		case r.Method == "PUT" && r.URL.Path == "/test-index":
			indexCreated = true
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"acknowledged":true}`))

		case r.Method == "POST" && r.URL.Path == "/test-index/_doc":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"result":"created"}`))

		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))
	defer server.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	if err := forwarder.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "new"})
	if err := forwarder.Call(t.Context(), record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !indexCreated {
		t.Errorf("expected index to be created")
	}

	if err := forwarder.Close(t.Context()); err != nil {
		t.Fatalf("failed to close forwarder: %v", err)
	}
}

func TestForwarderElastic_Call_IndexCreateFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")

		switch {
		case r.Method == "HEAD" && r.URL.Path == "/test-index":
			w.WriteHeader(http.StatusNotImplemented)

		case r.Method == "PUT" && r.URL.Path == "/test-index":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"fail"}`))

		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))
	defer server.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	if err := forwarder.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "fail"})
	if err := forwarder.Call(t.Context(), record); err == nil {
		t.Fatalf("expected error when index creation fails")
	}

	if err := forwarder.Close(t.Context()); err != nil {
		t.Fatalf("failed to close forwarder: %v", err)
	}
}

func TestForwarderElastic_Call_IndexFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")

		switch {
		case r.Method == "HEAD" && r.URL.Path == "/test-index":
			w.WriteHeader(http.StatusOK)

		case r.Method == "POST" && r.URL.Path == "/test-index/_doc":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"fail"}`))

		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))
	defer server.Close()

	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	if err := forwarder.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "fail"})
	if err := forwarder.Call(t.Context(), record); err == nil {
		t.Fatalf("expected error when indexing fails")
	}

	if err := forwarder.Close(t.Context()); err != nil {
		t.Fatalf("failed to close forwarder: %v", err)
	}
}

func TestForwarderElastic_Call_InvalidCACert(t *testing.T) {
	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{"http://localhost:9200"},
				CACert:    "not-pem",
			},
		},
	}

	if err := forwarder.Init(t.Context()); err == nil {
		t.Fatalf("expected error due to invalid CA cert")
	}
}
