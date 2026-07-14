package forwarders_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"link-society.com/flowg/internal/engines/forwarders"
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

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "hello"})
	if err := runtime.Call(t.Context(), record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := runtime.Close(t.Context()); err != nil {
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

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "new"})
	if err := runtime.Call(t.Context(), record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !indexCreated {
		t.Errorf("expected index to be created")
	}

	if err := runtime.Close(t.Context()); err != nil {
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

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "fail"})
	if err := runtime.Call(t.Context(), record); err == nil {
		t.Fatalf("expected error when index creation fails")
	}

	if err := runtime.Close(t.Context()); err != nil {
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

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err != nil {
		t.Fatalf("failed to initialize forwarder: %v", err)
	}

	record := models.NewLogRecord(map[string]string{"msg": "fail"})
	if err := runtime.Call(t.Context(), record); err == nil {
		t.Fatalf("expected error when indexing fails")
	}

	if err := runtime.Close(t.Context()); err != nil {
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

	runtime, err := forwarders.NewRuntime(forwarder)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	if err := runtime.Init(t.Context()); err == nil {
		t.Fatalf("expected error due to invalid CA cert")
	}
}
