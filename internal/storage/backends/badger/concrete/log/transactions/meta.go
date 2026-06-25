package transactions

import (
	"fmt"

	"encoding/json"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/app/featureflags"

	"link-society.com/flowg/internal/models"
)

var demoStreamConfig = models.StreamConfig{
	RetentionTime: 15 * 60, // 15 minutes
	RetentionSize: 10,
	IndexedFields: []string{},
}

// FetchStreamConfigs returns the configuration of every stream, keyed by stream
// name. In demo mode the fixed demo configuration is substituted for all of
// them.
func FetchStreamConfigs(txn *badger.Txn) (map[string]models.StreamConfig, error) {
	streams := map[string]models.StreamConfig{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = []byte("stream:config:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		stream := string(it.Item().Key()[len(opts.Prefix):])

		var streamConfig models.StreamConfig
		if featureflags.GetDemoMode() {
			streamConfig = demoStreamConfig
		} else {
			err := it.Item().Value(func(val []byte) error {
				if len(val) > 0 {
					if err := json.Unmarshal(val, &streamConfig); err != nil {
						return fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}

		if streamConfig.IndexedFields == nil {
			streamConfig.IndexedFields = []string{}
		}

		streams[stream] = streamConfig
	}

	return streams, nil
}

// FetchStreamFields lists the field names recorded for a stream from its
// "stream:field:<stream>:" existence markers.
func FetchStreamFields(txn *badger.Txn, stream string) []string {
	fields := []string{}

	prefix := fmt.Appendf(nil, "stream:field:%s:", stream)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = prefix
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		key := item.Key()

		fieldName := string(key[len(prefix):])
		fields = append(fields, fieldName)
	}

	return fields
}

// GetOrCreateStreamConfig returns a stream's configuration, lazily creating an
// empty (default) config the first time a stream is referenced so later writes
// have something to update. Demo mode short-circuits to the fixed config.
func GetOrCreateStreamConfig(txn *badger.Txn, stream string) (models.StreamConfig, error) {
	if featureflags.GetDemoMode() {
		return demoStreamConfig, nil
	}

	var streamConfig models.StreamConfig

	streamKey := fmt.Appendf(nil, "stream:config:%s", stream)
	switch streamConfigItem, err := txn.Get(streamKey); {
	case err != nil && err != badger.ErrKeyNotFound:
		return models.StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)

	case err == badger.ErrKeyNotFound:
		err := txn.Set(streamKey, []byte(""))
		if err != nil {
			return models.StreamConfig{}, fmt.Errorf(
				"could not create default stream config '%s': %w",
				stream, err,
			)
		}

	case err == nil:
		err := streamConfigItem.Value(func(val []byte) error {
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
					return fmt.Errorf(
						"could not unmarshal stream config '%s': %w",
						stream, err,
					)
				}
			}
			return nil
		})
		if err != nil {
			return models.StreamConfig{}, err
		}
	}

	if streamConfig.IndexedFields == nil {
		streamConfig.IndexedFields = []string{}
	}

	return streamConfig, nil
}

// ConfigureStream stores a new stream configuration and reconciles its indexed
// fields against the previous one: fields newly added to the index are back-
// filled over existing records, and fields removed from it have their index
// dropped.
func ConfigureStream(txn *badger.Txn, stream string, config models.StreamConfig) error {
	if config.IndexedFields == nil {
		config.IndexedFields = []string{}
	}

	oldConfig, err := GetOrCreateStreamConfig(txn, stream)
	if err != nil {
		return fmt.Errorf("could not fetch old stream config '%s': %w", stream, err)
	}

	streamKey := fmt.Appendf(nil, "stream:config:%s", stream)
	configVal, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
	}

	if err := txn.Set(streamKey, configVal); err != nil {
		return fmt.Errorf("could not save stream config '%s': %w", stream, err)
	}

	for _, field := range config.IndexedFields {
		if !oldConfig.IsFieldIndexed(field) {
			IndexField(txn, stream, field)
		}
	}

	for _, field := range oldConfig.IndexedFields {
		if !config.IsFieldIndexed(field) {
			UnindexField(txn, stream, field)
		}
	}

	return nil
}

// DeleteStream removes everything belonging to a stream: its log entries
// ("entry:<stream>:*"), inverted indexes ("index:<stream>:*"), field markers
// ("stream:field:<stream>:*") and its configuration key.
func DeleteStream(txn *badger.Txn, stream string) error {
	prefixes := []string{
		fmt.Sprintf("entry:%s:", stream),
		fmt.Sprintf("index:%s:", stream),
		fmt.Sprintf("stream:field:%s:", stream),
	}

	for _, prefix := range prefixes {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)

			if err := txn.Delete(key); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from stream '%s': %w",
					key, stream, err,
				)
			}
		}
	}

	streamKey := fmt.Appendf(nil, "stream:config:%s", stream)
	if err := txn.Delete(streamKey); err != nil {
		return fmt.Errorf("could not delete stream config '%s': %w", stream, err)
	}

	return nil
}
