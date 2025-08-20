package middlewares_test

import (
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"context"

	"bytes"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/api/middlewares"
)

func TestElasticEndpoint(t *testing.T) {
	mockAuthStorage := auth.NewMockStorage().(*auth.MockStorage)
	mockLogStorage := log.NewMockStorage().(*log.MockStorage)
	mockConfigStorage := config.NewMockStorage().(*config.MockStorage)

	mockLogNotifier := lognotify.NewMockNotifier().(*lognotify.MockNotifier)
	mockPipelineRunner := pipelines.NewMockRunner().(*pipelines.MockRunner)

	deps := &middlewares.Dependencies{
		AuthStorage:   mockAuthStorage,
		LogStorage:    mockLogStorage,
		ConfigStorage: mockConfigStorage,

		LogNotifier:    mockLogNotifier,
		PipelineRunner: mockPipelineRunner,
	}

	mockAuthStorage.On("VerifyUserPassword", mock.Anything, "test", "test").
		Return(true, nil)
	mockAuthStorage.On("FetchUser", mock.Anything, "test").
		Return(
			&models.User{
				Name:  "test",
				Roles: []string{"admin"},
			},
			nil,
		)
	mockAuthStorage.On("VerifyUserPermission", mock.Anything, "test", mock.Anything).
		Return(true, nil)

	mockConfigStorage.On("ListPipelines", mock.Anything).
		Return([]string{"test"}, nil)

	mockPipelineRunner.On("Run", mock.Anything, "test", pipelines.DIRECT_ENTRYPOINT, mock.Anything).
		Return(nil)

	handler := middlewares.NewHandler(deps)

	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := elasticsearch.Config{
		Username:  "test",
		Password:  "test",
		Addresses: []string{fmt.Sprintf("%s/api/v1/middlewares/elastic/", server.URL)},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create ElasticSearch client: %v", err)
	}

	ctx := context.Background()

	resp, err := client.Indices.Exists([]string{"test"}, client.Indices.Exists.WithContext(ctx))
	if err != nil {
		t.Fatalf("failed to send ElasticSearch request: %v", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		t.Fatalf("failed to check index: %s", resp.String())
	}

	data := bytes.NewReader([]byte(`{"message": "test log"}`))
	resp, err = client.Index("test", data, client.Index.WithContext(ctx))
	if err != nil {
		t.Fatalf("failed to send ElasticSearch request: %v", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		t.Fatalf("failed to index document: %s", resp.String())
	}

	mockAuthStorage.AssertExpectations(t)
	mockLogStorage.AssertExpectations(t)
	mockConfigStorage.AssertExpectations(t)

	mockLogNotifier.AssertExpectations(t)
	mockPipelineRunner.AssertExpectations(t)
}
