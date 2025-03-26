package models_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"context"
	"net/http"

	"maps"

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
				Url:     testServer.URL,
				ApiKey:  "apiKey",
				Source:  "sourceTag",
				Service: "serviceTag",
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
				Url:     testServer.URL,
				ApiKey:  "apiKey",
				Source:  "sourceTag",
				Service: "serviceTag",
			},
		},
	}

	record := models.NewLogRecord(map[string]string{})
	err := forwarder.Call(context.Background(), record)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestForwarderDatadog_CreateDatadogHttpLogItem(t *testing.T) {

	datadogForwarder := &models.ForwarderDatadogV2{
		Url:     "url",
		ApiKey:  "apiKey",
		Source:  "technology",
		Service: "service",
	}

	record := models.NewLogRecord(map[string]string{
		"technology": "nginx",
		"service":    "myservice",
		"hostname":   "localhost",
		"message":    "some log message",
	})

	logTags := map[string]string{
		"level": "INFO",
		"tag1":  "value1",
		"tag2":  "value2",
		"tag3":  "value3",
	}

	maps.Copy(record.Fields, logTags)

	datadogLogItem := models.CreateDatadogHttpLogItem(datadogForwarder, record)

	if datadogLogItem.Service != record.Fields[datadogForwarder.Service] {
		t.Errorf("Unexpected error computing datadog service. Got %s, want %s", datadogLogItem.Service, record.Fields[datadogForwarder.Service])
	}

	if datadogLogItem.DdSource != record.Fields[datadogForwarder.Source] {
		t.Errorf("Unexpected error computing datadog source. Got %s, want %s", datadogLogItem.DdSource, record.Fields[datadogForwarder.Source])
	}

	if datadogLogItem.Hostname != record.Fields["hostname"] {
		t.Errorf("Unexpected error computing datadog hostname. Got %s, want %s", datadogLogItem.Hostname, record.Fields["hostname"])
	}

	if datadogLogItem.Message != record.Fields["message"] {
		t.Errorf("Unexpected error computing datadog log message. Got %s, want %s", datadogLogItem.Message, record.Fields["message"])
	}

	for key, value := range logTags {
		tagString := key + ":" + value
		if !strings.Contains(datadogLogItem.DdTags, tagString) {
			t.Errorf("Unexpected error computing datadog tags. Expected %s to contain %s", datadogLogItem.DdTags, tagString)
		}
	}
}
