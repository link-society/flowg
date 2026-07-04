package transactions

import (
	"fmt"

	"encoding/json"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

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
func FetchStreamConfigs(txn fdb.ReadTransaction) (map[string]models.StreamConfig, error) {
	streams := map[string]models.StreamConfig{}

	ri := txn.GetRange(cfgSS, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()

		tpl, err := tuple.Unpack(kv.Key[len(cfgSS.Bytes()):])
		if err != nil {
			return nil, fmt.Errorf("could not unpack config key: %w", err)
		}
		if len(tpl) < 1 {
			continue
		}
		stream, ok := tpl[0].(string)
		if !ok {
			continue
		}

		var streamConfig models.StreamConfig
		if featureflags.GetDemoMode() {
			streamConfig = demoStreamConfig
		} else {
			if len(kv.Value) > 0 {
				if err := json.Unmarshal(kv.Value, &streamConfig); err != nil {
					return nil, fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
				}
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
// field subspace markers.
func FetchStreamFields(txn fdb.ReadTransaction, stream string) ([]string, error) {
	fields := []string{}

	streamFieldSS := fieldSS.Sub(subspace.FromBytes([]byte(stream)))

	ri := txn.GetRange(streamFieldSS, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()

		tpl, err := tuple.Unpack(kv.Key[len(streamFieldSS.Bytes()):])
		if err != nil {
			return nil, fmt.Errorf("could not unpack field key: %w", err)
		}
		if len(tpl) < 1 {
			continue
		}
		fieldName, ok := tpl[0].(string)
		if !ok {
			continue
		}

		fields = append(fields, fieldName)
	}

	return fields, nil
}

// GetOrCreateStreamConfig returns a stream's configuration, lazily creating an
// empty (default) config the first time a stream is referenced so later writes
// have something to update. Demo mode short-circuits to the fixed config.
func GetOrCreateStreamConfig(txn fdb.Transaction, stream string) (models.StreamConfig, error) {
	if featureflags.GetDemoMode() {
		return demoStreamConfig, nil
	}

	var streamConfig models.StreamConfig

	cfgKey := cfgSS.Pack(tuple.Tuple{stream})
	val, err := txn.Get(cfgKey).Get()
	if err != nil {
		return models.StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)
	}

	if val == nil {
		// Key does not exist — create default config.
		txn.Set(cfgKey, nil)
	} else if len(val) > 0 {
		if err := json.Unmarshal(val, &streamConfig); err != nil {
			return models.StreamConfig{}, fmt.Errorf(
				"could not unmarshal stream config '%s': %w",
				stream, err,
			)
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
func ConfigureStream(txn fdb.Transaction, stream string, config models.StreamConfig) error {
	if config.IndexedFields == nil {
		config.IndexedFields = []string{}
	}

	oldConfig, err := GetOrCreateStreamConfig(txn, stream)
	if err != nil {
		return fmt.Errorf("could not fetch old stream config '%s': %w", stream, err)
	}

	cfgKey := cfgSS.Pack(tuple.Tuple{stream})
	configVal, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
	}

	txn.Set(cfgKey, configVal)

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
// (entry sub-space), inverted indexes (index sub-space), field markers
// (field sub-space) and its configuration key.
func DeleteStream(txn fdb.Transaction, stream string) error {
	// Delete all entry keys for the stream.
	streamEntrySS := entrySS.Sub(subspace.FromBytes([]byte(stream)))
	txn.ClearRange(streamEntrySS)

	// Delete all index keys for the stream.
	streamIndexSS := indexSS.Sub(subspace.FromBytes([]byte(stream)))
	txn.ClearRange(streamIndexSS)

	// Delete all field markers for the stream.
	streamFieldSS := fieldSS.Sub(subspace.FromBytes([]byte(stream)))
	txn.ClearRange(streamFieldSS)

	// Delete the configuration key.
	cfgKey := cfgSS.Pack(tuple.Tuple{stream})
	txn.Clear(cfgKey)

	return nil
}
