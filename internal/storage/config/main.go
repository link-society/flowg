package config

import (
	"context"
	"fmt"

	"io"

	"encoding/json"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/config/transactions"
)

const (
	transformerItemType = "transformer"
	pipelineItemType    = "pipeline"
	forwarderItemType   = "forwarder"
)

type options struct {
	dir      string
	inMemory bool
	readOnly bool
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

func OptReadOnly(readOnly bool) func(*options) {
	return func(o *options) {
		o.readOnly = readOnly
	}
}

type Storage struct {
	proctree.Process

	kvStore *kvstore.Storage
}

var _ proctree.Process = (*Storage)(nil)

func NewStorage(opts ...func(*options)) *Storage {
	options := options{
		dir:      "./data/config",
		inMemory: false,
		readOnly: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	kvStore := kvstore.NewStorage(
		kvstore.OptLogChannel("configstorage"),
		kvstore.OptDirectory(options.dir),
		kvstore.OptInMemory(options.inMemory),
		kvstore.OptReadOnly(options.readOnly),
	)

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		kvStore,
		proctree.NewProcess(&migratorProcH{
			baseDir: options.dir,
			storage: &Storage{kvStore: kvStore},
		}),
	)

	return &Storage{
		Process: process,

		kvStore: kvStore,
	}
}

func (s *Storage) Backup(ctx context.Context, w io.Writer) error {
	return s.kvStore.Backup(ctx, w, 0)
}

func (s *Storage) Restore(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *Storage) ListTransformers(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, transformerItemType)
}

func (s *Storage) ReadTransformer(ctx context.Context, name string) (string, error) {
	content, err := s.readItem(ctx, transformerItemType, name)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *Storage) WriteTransformer(ctx context.Context, name string, content string) error {
	return s.writeItem(ctx, transformerItemType, name, []byte(content))
}

func (s *Storage) DeleteTransformer(ctx context.Context, name string) error {
	return s.deleteItem(ctx, transformerItemType, name)
}

func (s *Storage) ListPipelines(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, pipelineItemType)
}

func (s *Storage) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error) {
	content, err := s.readItem(ctx, pipelineItemType, name)
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

	return s.writeItem(ctx, pipelineItemType, name, content)
}

func (s *Storage) WriteRawPipeline(ctx context.Context, name string, content string) error {
	return s.writeItem(ctx, pipelineItemType, name, []byte(content))
}

func (s *Storage) DeletePipeline(ctx context.Context, name string) error {
	return s.deleteItem(ctx, pipelineItemType, name)
}

func (s *Storage) ListForwarders(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, forwarderItemType)
}

func (s *Storage) ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error) {
	content, err := s.readItem(ctx, forwarderItemType, name)
	if err != nil {
		return nil, err
	}

	webhook, changed, err := models.ConvertForwarder(content)
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

	return s.writeItem(ctx, forwarderItemType, name, content)
}

func (s *Storage) DeleteForwarder(ctx context.Context, name string) error {
	return s.deleteItem(ctx, forwarderItemType, name)
}

func (s *Storage) listItems(ctx context.Context, itemType string) ([]string, error) {
	var items []string

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			items, err = transactions.ListItems(txn, itemType)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Storage) readItem(
	ctx context.Context,
	itemType string,
	name string,
) ([]byte, error) {
	var content []byte

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			content, err = transactions.ReadItem(txn, itemType, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *Storage) writeItem(
	ctx context.Context,
	itemType string,
	name string,
	content []byte,
) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.WriteItem(txn, itemType, name, content)
		},
	)
}

func (s *Storage) deleteItem(ctx context.Context, itemType string, name string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteItem(txn, itemType, name)
		},
	)
}
