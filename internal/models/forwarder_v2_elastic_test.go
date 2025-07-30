package models_test

import (
	"net/http/httptest"
	"testing"

	"context"

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
		Config: &models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	record := models.NewLogRecord(map[string]string{"msg": "hello"})
	err := forwarder.Call(context.Background(), record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
		Config: &models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	record := models.NewLogRecord(map[string]string{"msg": "new"})
	err := forwarder.Call(context.Background(), record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !indexCreated {
		t.Errorf("expected index to be created")
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
		Config: &models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	record := models.NewLogRecord(map[string]string{"msg": "fail"})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error when index creation fails")
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
		Config: &models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{server.URL},
			},
		},
	}

	record := models.NewLogRecord(map[string]string{"msg": "fail"})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error when indexing fails")
	}
}

func TestForwarderElastic_Call_InvalidCACert(t *testing.T) {
	forwarder := &models.ForwarderV2{
		Version: 2,
		Config: &models.ForwarderConfigV2{
			Elastic: &models.ForwarderElasticV2{
				Type:      "elastic",
				Index:     "test-index",
				Addresses: []string{"http://localhost:9200"},
				CACert:    "not-base64",
			},
		},
	}

	record := models.NewLogRecord(map[string]string{"msg": "bad ca"})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error due to invalid CA cert")
	}
}
