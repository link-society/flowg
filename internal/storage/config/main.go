package config

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"io"

	"encoding/json"

	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/config/transactions"
	"link-society.com/flowg/internal/utils/kvstore"
)

type Storage interface {
	storage.Streamable

	ListTransformers(ctx context.Context) ([]string, error)
	ReadTransformer(ctx context.Context, name string) (string, error)
	WriteTransformer(ctx context.Context, name string, content string) error
	DeleteTransformer(ctx context.Context, name string) error

	ListPipelines(ctx context.Context) ([]string, error)
	ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error)
	WritePipeline(ctx context.Context, name string, flow *models.FlowGraphV2) error
	WriteRawPipeline(ctx context.Context, name string, content string) error
	DeletePipeline(ctx context.Context, name string) error

	ListForwarders(ctx context.Context) ([]string, error)
	ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error)
	WriteForwarder(ctx context.Context, name string, forwarder *models.ForwarderV2) error
	DeleteForwarder(ctx context.Context, name string) error

	ReadSystemConfig(ctx context.Context) (*models.SystemConfiguration, error)
	WriteSystemConfig(ctx context.Context, config *models.SystemConfiguration) error
}

const (
	transformerItemType = "transformer"
	pipelineItemType    = "pipeline"
	forwarderItemType   = "forwarder"
	systemItemType      = "system"
)

type Options struct {
	Directory string
	InMemory  bool
	ReadOnly  bool
}

type storageImpl struct {
	kvStore               kvstore.Storage
	lock                  *sync.Mutex
	configurationInstance *models.SystemConfiguration
}

type deps struct {
	fx.In

	S kvstore.Storage `name:"storage.config"`
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		Directory: "",
		InMemory:  false,
		ReadOnly:  false,
	}
}

func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "storage.config"
	kvOpts.Directory = opts.Directory
	kvOpts.InMemory = opts.InMemory
	kvOpts.ReadOnly = opts.ReadOnly

	return fx.Module(
		"storage.config",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) Storage {
			storage := &storageImpl{
				kvStore:               d.S,
				lock:                  &sync.Mutex{},
				configurationInstance: nil,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := migrateAlerts(opts.Directory); err != nil {
						return fmt.Errorf("failed to migrate alerts: %w", err)
					}

					if err := migrateToBadger(ctx, opts.Directory, storage); err != nil {
						return fmt.Errorf("failed to migrate to badger: %w", err)
					}

					return nil
				},
			})

			return storage
		}),
	)
}

func (s *storageImpl) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return s.kvStore.Backup(ctx, w, since)
}

func (s *storageImpl) Load(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *storageImpl) ListTransformers(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, transformerItemType)
}

func (s *storageImpl) ReadTransformer(ctx context.Context, name string) (string, error) {
	content, err := s.readItem(ctx, transformerItemType, name)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *storageImpl) WriteTransformer(ctx context.Context, name string, content string) error {
	return s.writeItem(ctx, transformerItemType, name, []byte(content))
}

func (s *storageImpl) DeleteTransformer(ctx context.Context, name string) error {
	return s.deleteItem(ctx, transformerItemType, name)
}

func (s *storageImpl) ListPipelines(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, pipelineItemType)
}

func (s *storageImpl) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error) {
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

func (s *storageImpl) WritePipeline(ctx context.Context, name string, flow *models.FlowGraphV2) error {
	content, err := json.Marshal(flow)
	if err != nil {
		return fmt.Errorf("failed to marshal flow: %w", err)
	}

	return s.writeItem(ctx, pipelineItemType, name, content)
}

func (s *storageImpl) WriteRawPipeline(ctx context.Context, name string, content string) error {
	return s.writeItem(ctx, pipelineItemType, name, []byte(content))
}

func (s *storageImpl) DeletePipeline(ctx context.Context, name string) error {
	return s.deleteItem(ctx, pipelineItemType, name)
}

func (s *storageImpl) ListForwarders(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, forwarderItemType)
}

func (s *storageImpl) ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error) {
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

func (s *storageImpl) WriteForwarder(ctx context.Context, name string, forwarder *models.ForwarderV2) error {
	content, err := json.Marshal(forwarder)
	if err != nil {
		return fmt.Errorf("failed to marshal forwarder: %w", err)
	}

	return s.writeItem(ctx, forwarderItemType, name, content)
}

func (s *storageImpl) DeleteForwarder(ctx context.Context, name string) error {
	return s.deleteItem(ctx, forwarderItemType, name)
}

func (s *storageImpl) ReadSystemConfig(ctx context.Context) (*models.SystemConfiguration, error) {
	if s.configurationInstance != nil {
		return s.configurationInstance, nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.configurationInstance == nil {
		content, err := s.readItem(ctx, systemItemType, "config")
		if errors.Is(err, badger.ErrKeyNotFound) {
			s.configurationInstance = &models.SystemConfiguration{}
			return s.configurationInstance, nil
		}

		if err != nil {
			return nil, err
		}

		var config models.SystemConfiguration
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, err
		}

		s.configurationInstance = &config
	}

	return s.configurationInstance, nil
}

func (s *storageImpl) WriteSystemConfig(ctx context.Context, config *models.SystemConfiguration) error {
	if config.SyslogAllowedOrigins != nil {
		for _, origin := range config.SyslogAllowedOrigins {
			if strings.Contains(origin, "/") {
				_, _, err := net.ParseCIDR(origin)
				if err != nil {
					return fmt.Errorf("invalid syslog allow origin: %s", origin)
				}
			} else {
				if net.ParseIP(origin) == nil {
					return fmt.Errorf("invalid syslog allow origin: %s", origin)
				}
			}
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.configurationInstance = nil

	content, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return s.writeItem(ctx, systemItemType, "config", content)
}

func (s *storageImpl) listItems(ctx context.Context, itemType string) ([]string, error) {
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

func (s *storageImpl) readItem(
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

func (s *storageImpl) writeItem(
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

func (s *storageImpl) deleteItem(ctx context.Context, itemType string, name string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteItem(txn, itemType, name)
		},
	)
}
