package config

import (
	"context"
	"errors"
	"fmt"

	"encoding/base64"
	"encoding/json"
	"path/filepath"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/filestore"
)

type options struct {
	dir      string
	inMemory bool
}

func OptDirectory(dir string) func(*options) {
	return func(o *options) {
		o.dir = dir
	}
}

func OptInMemory(inMemory bool) func(*options) {
	return func(o *options) {
		o.inMemory = inMemory
	}
}

type Storage struct {
	transformerStore *filestore.Storage
	pipelineStore    *filestore.Storage
	alertStore       *filestore.Storage

	actor actor.Actor
}

func NewStorage(opts ...func(*options)) *Storage {
	options := options{
		dir:      "./data/config",
		inMemory: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	transformerStore := filestore.NewStorage(
		filestore.OptDirectory(filepath.Join(options.dir, "transformers")),
		filestore.OptInMemory(options.inMemory),
		filestore.OptExtension(".vrl"),
	)
	pipelineStore := filestore.NewStorage(
		filestore.OptDirectory(filepath.Join(options.dir, "pipelines")),
		filestore.OptInMemory(options.inMemory),
		filestore.OptExtension(".json"),
	)
	alertStore := filestore.NewStorage(
		filestore.OptDirectory(filepath.Join(options.dir, "alerts")),
		filestore.OptInMemory(options.inMemory),
		filestore.OptExtension(".json.b64"),
	)

	actor := actor.Combine(transformerStore, pipelineStore, alertStore).
		WithOptions(actor.OptStopTogether()).
		Build()

	return &Storage{
		transformerStore: transformerStore,
		pipelineStore:    pipelineStore,
		alertStore:       alertStore,

		actor: actor,
	}
}

func (s *Storage) Start() {
	s.actor.Start()
}

func (s *Storage) WaitStarted() error {
	errs := []error{}

	if err := s.transformerStore.WaitStarted(); err != nil {
		errs = append(errs, fmt.Errorf("failed to start transformer store: %w", err))
	}

	if err := s.pipelineStore.WaitStarted(); err != nil {
		errs = append(errs, fmt.Errorf("failed to start pipeline store: %w", err))
	}

	if err := s.alertStore.WaitStarted(); err != nil {
		errs = append(errs, fmt.Errorf("failed to start alert store: %w", err))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (s *Storage) Stop() {
	s.actor.Stop()
}

func (s *Storage) WaitStopped() error {
	errs := []error{}

	if err := s.alertStore.WaitStopped(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop alert store: %w", err))
	}

	if err := s.pipelineStore.WaitStopped(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop pipeline store: %w", err))
	}

	if err := s.transformerStore.WaitStopped(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop transformer store: %w", err))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (s *Storage) ListTransformers(ctx context.Context) ([]string, error) {
	return s.transformerStore.ListFiles(ctx)
}

func (s *Storage) ReadTransformer(ctx context.Context, name string) (string, error) {
	content, err := s.transformerStore.ReadFile(ctx, name)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *Storage) WriteTransformer(ctx context.Context, name string, content string) error {
	return s.transformerStore.WriteFile(ctx, name, []byte(content))
}

func (s *Storage) DeleteTransformer(ctx context.Context, name string) error {
	return s.transformerStore.DeleteFile(ctx, name)
}

func (s *Storage) ListPipelines(ctx context.Context) ([]string, error) {
	return s.pipelineStore.ListFiles(ctx)
}

func (s *Storage) ReadPipeline(ctx context.Context, name string) (*models.FlowGraph, error) {
	content, err := s.pipelineStore.ReadFile(ctx, name)
	if err != nil {
		return nil, err
	}

	var flowGraph *models.FlowGraph
	if err := json.Unmarshal(content, &flowGraph); err != nil {
		return nil, fmt.Errorf("failed to unmarshal flow: %w", err)
	}

	return flowGraph, nil
}

func (s *Storage) WritePipeline(ctx context.Context, name string, flow *models.FlowGraph) error {
	content, err := json.Marshal(flow)
	if err != nil {
		return fmt.Errorf("failed to marshal flow: %w", err)
	}

	return s.pipelineStore.WriteFile(ctx, name, content)
}

func (s *Storage) WriteRawPipeline(ctx context.Context, name string, content string) error {
	return s.pipelineStore.WriteFile(ctx, name, []byte(content))
}

func (s *Storage) DeletePipeline(ctx context.Context, name string) error {
	return s.pipelineStore.DeleteFile(ctx, name)
}

func (s *Storage) ListAlerts(ctx context.Context) ([]string, error) {
	return s.alertStore.ListFiles(ctx)
}

func (s *Storage) ReadAlert(ctx context.Context, name string) (*models.Webhook, error) {
	b64content, err := s.alertStore.ReadFile(ctx, name)
	if err != nil {
		return nil, err
	}

	content := make([]byte, base64.StdEncoding.DecodedLen(len(b64content)))
	n, err := base64.StdEncoding.Decode(content, b64content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode webhook %s: %w", name, err)
	}

	var webhook *models.Webhook
	if err := json.Unmarshal(content[:n], &webhook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook: %w", err)
	}

	return webhook, nil
}

func (s *Storage) WriteAlert(ctx context.Context, name string, webhook *models.Webhook) error {
	content, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook: %w", err)
	}

	b64content := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
	base64.StdEncoding.Encode(b64content, content)

	return s.alertStore.WriteFile(ctx, name, b64content)
}

func (s *Storage) DeleteAlert(ctx context.Context, name string) error {
	return s.alertStore.DeleteFile(ctx, name)
}
