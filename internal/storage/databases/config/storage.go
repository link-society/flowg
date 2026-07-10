package config

import (
	"context"
	"fmt"

	"sync"

	"encoding/json"
	"strings"

	"io"
	"net"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/storage/databases/config/transactions"
)

const (
	transformerItemType = "transformer"
	pipelineItemType    = "pipeline"
	forwarderItemType   = "forwarder"
	systemItemType      = "system"
)

// Storage is a backend-agnostic implementation of [storage.ConfigStorage]. It
// runs the config transactions from the transactions subpackage on top of any
// [kv.Adapter] and caches the decoded system configuration in memory.
type Storage[QTx kv.QueryTx, MTx kv.MutationTx] struct {
	adapter kv.Adapter[QTx, MTx]

	lock                  *sync.Mutex
	configurationInstance *models.SystemConfiguration
}

var _ storage.ConfigStorage = (*Storage[kv.QueryTx, kv.MutationTx])(nil)

// NewStorage returns a [Storage] that persists configuration data through the
// given key-value adapter.
func NewStorage[QTx kv.QueryTx, MTx kv.MutationTx](adapter kv.Adapter[QTx, MTx]) *Storage[QTx, MTx] {
	return &Storage[QTx, MTx]{
		adapter:               adapter,
		lock:                  &sync.Mutex{},
		configurationInstance: nil,
	}
}

// Dump implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) Dump(ctx context.Context, w io.Writer, version uint64) (uint64, error) {
	return s.adapter.Backup(ctx, w, version)
}

// Load implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) Load(ctx context.Context, r io.Reader) error {
	return s.adapter.Restore(ctx, r)
}

// ListTransformers implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ListTransformers(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, transformerItemType)
}

// ReadTransformer implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ReadTransformer(ctx context.Context, name string) (string, error) {
	content, err := s.readItem(ctx, transformerItemType, name)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// WriteTransformer implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) WriteTransformer(ctx context.Context, name string, content string) error {
	return s.writeItem(ctx, transformerItemType, name, []byte(content))
}

// DeleteTransformer implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) DeleteTransformer(ctx context.Context, name string) error {
	return s.deleteItem(ctx, transformerItemType, name)
}

// ListPipelines implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ListPipelines(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, pipelineItemType)
}

// ReadPipeline implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error) {
	content, err := s.readItem(ctx, pipelineItemType, name)
	if err != nil {
		return nil, err
	}

	flowGraph, changed, err := models.ConvertFlowGraph(content)
	if err != nil {
		return nil, err
	}

	if changed {
		// If the flow graph was converted to a new version, we should save it back to the storage.
		if err := s.WritePipeline(ctx, name, flowGraph); err != nil {
			return nil, fmt.Errorf("failed to write updated flow graph: %w", err)
		}
	}

	return flowGraph, nil
}

// WritePipeline implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) WritePipeline(ctx context.Context, name string, flow *models.FlowGraphV2) error {
	content, err := json.Marshal(flow)
	if err != nil {
		return fmt.Errorf("failed to marshal flow graph: %w", err)
	}

	return s.writeItem(ctx, pipelineItemType, name, content)
}

// WriteRawPipeline implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) WriteRawPipeline(ctx context.Context, name string, content string) error {
	return s.writeItem(ctx, pipelineItemType, name, []byte(content))
}

// DeletePipeline implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) DeletePipeline(ctx context.Context, name string) error {
	return s.deleteItem(ctx, pipelineItemType, name)
}

// ListForwarders implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ListForwarders(ctx context.Context) ([]string, error) {
	return s.listItems(ctx, forwarderItemType)
}

// ReadForwarder implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error) {
	content, err := s.readItem(ctx, forwarderItemType, name)
	if err != nil {
		return nil, err
	}

	webhook, changed, err := models.ConvertForwarder(content)
	if err != nil {
		return nil, err
	}

	if changed {
		// If the forwarder was converted to a new version, we should save it back
		// to the storage.
		if err := s.WriteForwarder(ctx, name, webhook); err != nil {
			return nil, fmt.Errorf("failed to write updated forwarder: %w", err)
		}
	}

	return webhook, nil
}

// WriteForwarder implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) WriteForwarder(ctx context.Context, name string, forwarder *models.ForwarderV2) error {
	content, err := json.Marshal(forwarder)
	if err != nil {
		return fmt.Errorf("failed to marshal forwarder: %w", err)
	}

	return s.writeItem(ctx, forwarderItemType, name, content)
}

// DeleteForwarder implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) DeleteForwarder(ctx context.Context, name string) error {
	return s.deleteItem(ctx, forwarderItemType, name)
}

// HasSystemConfig implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) HasSystemConfig(ctx context.Context) (bool, error) {
	val, err := s.readItem(ctx, systemItemType, "config")
	if err != nil {
		return false, err
	}
	return val != nil, nil
}

// ReadSystemConfig implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) ReadSystemConfig(ctx context.Context) (*models.SystemConfiguration, error) {
	if s.configurationInstance != nil {
		return s.configurationInstance, nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.configurationInstance == nil {
		content, err := s.readItem(ctx, systemItemType, "config")
		if err != nil {
			return nil, err
		}

		if content == nil {
			s.configurationInstance = &models.SystemConfiguration{}
			return s.configurationInstance, nil
		}

		var config models.SystemConfiguration
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal system configuration: %w", err)
		}

		s.configurationInstance = &config
	}

	return s.configurationInstance, nil
}

// WriteSystemConfig implements [storage.ConfigStorage].
func (s *Storage[QTx, MTx]) WriteSystemConfig(ctx context.Context, config *models.SystemConfiguration) error {
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

func (s *Storage[QTx, MTx]) listItems(ctx context.Context, itemType string) ([]string, error) {
	var items []string

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		items, err = transactions.ListItems(txn, itemType)
		return err
	})

	return items, err
}

func (s *Storage[QTx, MTx]) readItem(ctx context.Context, itemType string, name string) ([]byte, error) {
	var content []byte

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		content, err = transactions.ReadItem(txn, itemType, name)
		return err
	})

	return content, err
}

func (s *Storage[QTx, MTx]) writeItem(ctx context.Context, itemType string, name string, content []byte) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.WriteItem(txn, itemType, name, content)
	})
}

func (s *Storage[QTx, MTx]) deleteItem(ctx context.Context, itemType string, name string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.DeleteItem(txn, itemType, name)
	})
}
