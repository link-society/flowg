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

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/filestore"
	"link-society.com/flowg/internal/utils/proctree"
)

const (
	transformerExt = ".vrl"
	pipelineExt    = ".json"
	forwarderExt   = ".json.b64"
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
	proctree.Process

	transformerStore *filestore.Storage
	pipelineStore    *filestore.Storage
	forwarderStore   *filestore.Storage
}

var _ proctree.Process = (*Storage)(nil)

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
	forwarderStore := filestore.NewStorage(
		filestore.OptDirectory(filepath.Join(options.dir, "forwarders")),
		filestore.OptInMemory(options.inMemory),
		filestore.OptExtension(forwarderExt),
	)

	children := []proctree.Process{
		transformerStore,
		pipelineStore,
		forwarderStore,
	}

	if !options.inMemory {
		migratorPH := proctree.NewProcess(&migratorProcH{baseDir: options.dir})
		children = append([]proctree.Process{migratorPH}, children...)
	}

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		children...,
	)

	return &Storage{
		Process: process,

		transformerStore: transformerStore,
		pipelineStore:    pipelineStore,
		forwarderStore:   forwarderStore,
	}
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
		{storage: s.forwarderStore, kind: "forwarder"},
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

			case "forwarders":
				storage = s.forwarderStore
				ext = forwarderExt

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

func (s *Storage) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error) {
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

func (s *Storage) WritePipeline(ctx context.Context, name string, flow *models.FlowGraphV2) error {
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

func (s *Storage) ListForwarders(ctx context.Context) ([]string, error) {
	return s.forwarderStore.ListFiles(ctx)
}

func (s *Storage) ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error) {
	b64content, err := s.forwarderStore.ReadFile(ctx, name)
	if err != nil {
		return nil, err
	}

	content := make([]byte, base64.StdEncoding.DecodedLen(len(b64content)))
	n, err := base64.StdEncoding.Decode(content, b64content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode forwarder %s: %w", name, err)
	}

	webhook, changed, err := models.ConvertForwarder(content[:n])
	if err != nil {
		return nil, err
	}

	if changed {
		if err := s.WriteForwarder(ctx, name, webhook); err != nil {
			return nil, fmt.Errorf("failed to write updated forwarder: %w", err)
		}
	}

	return webhook, nil
}

func (s *Storage) WriteForwarder(ctx context.Context, name string, forwarder *models.ForwarderV2) error {
	content, err := json.Marshal(forwarder)
	if err != nil {
		return fmt.Errorf("failed to marshal forwarder: %w", err)
	}

	b64content := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
	base64.StdEncoding.Encode(b64content, content)

	return s.forwarderStore.WriteFile(ctx, name, b64content)
}

func (s *Storage) DeleteForwarder(ctx context.Context, name string) error {
	return s.forwarderStore.DeleteFile(ctx, name)
}
