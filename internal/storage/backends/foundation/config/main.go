package config

import (
	"context"
	"errors"
	"fmt"

	"encoding/json"
	"strings"

	"io"
	"net"

	"sync"

	"go.uber.org/fx"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/backends/foundation/config/transactions"
	"link-society.com/flowg/internal/storage/backends/foundation/kvstore"
)

var ErrNotSupported = errors.New("operation not supported")

const (
	transformerItemType = "transformer"
	pipelineItemType    = "pipeline"
	forwarderItemType   = "forwarder"
	systemItemType      = "system"
)

type Options struct {
	// ClusterFile is the path to the FoundationDB cluster file, which contains
	// the connection information for the database cluster.
	ClusterFile string

	// KeySpace is the key space to use for this database. All keys will be
	// prefixed with this key space, allowing multiple logical databases to share
	// the same physical cluster.
	KeySpace string
}

type storageImpl struct {
	kvStore               kvstore.Storage
	keySpace              subspace.Subspace
	lock                  *sync.Mutex
	configurationInstance *models.SystemConfiguration
}

type deps struct {
	fx.In

	S kvstore.Storage `name:"storage.config"`
}

var _ storage.ConfigStorage = (*storageImpl)(nil)

// DefaultOptions returns the default [Options] for the configuration storage.
func DefaultOptions() Options {
	return Options{
		ClusterFile: "",
		KeySpace:    "flowg",
	}
}

// NewStorage returns an fx module that provides a [storage.ConfigStorage]
// backed by FoundationDB.
func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.Tag = "config"
	kvOpts.ClusterFile = opts.ClusterFile

	return fx.Module(
		"storage.config",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(d deps) storage.ConfigStorage {
			return &storageImpl{
				kvStore:               d.S,
				keySpace:              subspace.FromBytes(fmt.Appendf(nil, "%s/config/", opts.KeySpace)),
				lock:                  &sync.Mutex{},
				configurationInstance: nil,
			}
		}),
	)
}

func (s *storageImpl) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return 0, ErrNotSupported
}

func (s *storageImpl) Load(ctx context.Context, r io.Reader) error {
	return ErrNotSupported
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

func (s *storageImpl) HasSystemConfig(ctx context.Context) (bool, error) {
	val, err := s.readItem(ctx, systemItemType, "config")
	if val == nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *storageImpl) ReadSystemConfig(ctx context.Context) (*models.SystemConfiguration, error) {
	if s.configurationInstance != nil {
		return s.configurationInstance, nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.configurationInstance == nil {
		content, err := s.readItem(ctx, systemItemType, "config")
		if content == nil {
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

	err := s.kvStore.View(ctx, func(txn fdb.ReadTransaction) error {
		var err error
		items, err = transactions.ListItems(txn, s.keySpace, itemType)
		return err
	})

	return items, err
}

func (s *storageImpl) readItem(ctx context.Context, itemType string, name string) ([]byte, error) {
	var content []byte

	err := s.kvStore.View(ctx, func(txn fdb.ReadTransaction) error {
		var err error
		content, err = transactions.ReadItem(txn, s.keySpace, itemType, name)
		return err
	})

	return content, err
}

func (s *storageImpl) writeItem(ctx context.Context, itemType string, name string, content []byte) error {
	return s.kvStore.Update(ctx, func(txn fdb.Transaction) error {
		return transactions.WriteItem(txn, s.keySpace, itemType, name, content)
	})
}

func (s *storageImpl) deleteItem(ctx context.Context, itemType string, name string) error {
	return s.kvStore.Update(ctx, func(txn fdb.Transaction) error {
		return transactions.DeleteItem(txn, s.keySpace, itemType, name)
	})
}
