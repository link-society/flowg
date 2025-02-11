package config

import (
	"context"
	"errors"
	"fmt"

	"archive/tar"
	"compress/gzip"
	"io"

	"encoding/base64"
	"encoding/json"
	"path/filepath"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/filestore"
)

const (
	transformerExt = ".vrl"
	pipelineExt    = ".json"
	alertExt       = ".json.b64"
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
		filestore.OptExtension(transformerExt),
	)
	pipelineStore := filestore.NewStorage(
		filestore.OptDirectory(filepath.Join(options.dir, "pipelines")),
		filestore.OptInMemory(options.inMemory),
		filestore.OptExtension(pipelineExt),
	)
	alertStore := filestore.NewStorage(
		filestore.OptDirectory(filepath.Join(options.dir, "alerts")),
		filestore.OptInMemory(options.inMemory),
		filestore.OptExtension(alertExt),
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

func (s *Storage) Backup(ctx context.Context, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	stores := []struct {
		storage *filestore.Storage
		kind    string
	}{
		{storage: s.transformerStore, kind: "transformer"},
		{storage: s.pipelineStore, kind: "pipeline"},
		{storage: s.alertStore, kind: "alert"},
	}

	for _, store := range stores {
		items, err := store.storage.ListFiles(ctx)
		if err != nil {
			return fmt.Errorf("failed to list %ss: %w", store.kind, err)
		}

		for _, name := range items {
			info, err := store.storage.StatFile(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to stat %s %s: %w", store.kind, name, err)
			}

			content, err := store.storage.ReadFile(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to read %s %s: %w", store.kind, name, err)
			}

			hdr, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return fmt.Errorf(
					"failed to create TAR header for %s %s: %w",
					store.kind,
					name,
					err,
				)
			}
			hdr.Name = filepath.Join(fmt.Sprintf("%ss", store.kind), info.Name())

			err = tw.WriteHeader(hdr)
			if err != nil {
				return fmt.Errorf(
					"failed to write TAR header for %s %s: %w",
					store.kind,
					name,
					err,
				)
			}

			_, err = tw.Write(content)
			if err != nil {
				return fmt.Errorf(
					"failed to write TAR content for %s %s: %w",
					store.kind,
					name,
					err,
				)
			}
		}
	}

	return nil
}

func (s *Storage) Restore(ctx context.Context, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create GZIP reader: %w", err)
	}

	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		switch hdr.Typeflag {
		case tar.TypeReg:
			kind := filepath.Base(filepath.Dir(hdr.Name))

			var (
				storage *filestore.Storage
				ext     string
			)
			switch kind {
			case "transformers":
				storage = s.transformerStore
				ext = transformerExt

			case "pipelines":
				storage = s.pipelineStore
				ext = pipelineExt

			case "alerts":
				storage = s.alertStore
				ext = alertExt

			default:
				return fmt.Errorf("unknown configuration item kind %s", kind)
			}

			data := make([]byte, hdr.Size)
			_, err := io.ReadFull(tr, data)
			if err != nil {
				return fmt.Errorf("failed to read TAR content: %w", err)
			}

			name := filepath.Base(hdr.Name)
			name = name[:len(name)-len(ext)]

			err = storage.WriteFile(ctx, name, data)
			if err != nil {
				return fmt.Errorf("failed to write %s %s: %w", kind, name, err)
			}
		}
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

func (s *Storage) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV1, error) {
	content, err := s.pipelineStore.ReadFile(ctx, name)
	if err != nil {
		return nil, err
	}

	flowGraph, changed, err := models.ConvertFlowGraph(content)
	if err != nil {
		return nil, err
	}

	if changed {
		if err := s.WritePipeline(ctx, name, flowGraph); err != nil {
			return nil, fmt.Errorf("failed to write updated flow graph: %w", err)
		}
	}

	return flowGraph, nil
}

func (s *Storage) WritePipeline(ctx context.Context, name string, flow *models.FlowGraphV1) error {
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

func (s *Storage) ReadAlert(ctx context.Context, name string) (*models.WebhookV1, error) {
	b64content, err := s.alertStore.ReadFile(ctx, name)
	if err != nil {
		return nil, err
	}

	content := make([]byte, base64.StdEncoding.DecodedLen(len(b64content)))
	n, err := base64.StdEncoding.Decode(content, b64content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode webhook %s: %w", name, err)
	}

	webhook, changed, err := models.ConvertWebhook(content[:n])
	if err != nil {
		return nil, err
	}

	if changed {
		if err := s.WriteAlert(ctx, name, webhook); err != nil {
			return nil, fmt.Errorf("failed to write updated webhook: %w", err)
		}
	}

	return webhook, nil
}

func (s *Storage) WriteAlert(ctx context.Context, name string, webhook *models.WebhookV1) error {
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
