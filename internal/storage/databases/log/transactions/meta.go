package transactions

import (
	"fmt"

	"encoding/json"

	"link-society.com/flowg/internal/app/featureflags"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

var demoStreamConfig = models.StreamConfig{
	RetentionTime: 15 * 60, // 15 minutes
	RetentionSize: 10,
	IndexedFields: []string{},
}

// FetchStreamConfigs returns the configuration of every stream, keyed by stream
// name. In demo mode the fixed demo configuration is substituted for all of
// them.
func FetchStreamConfigs(txn kv.QueryTx) (map[string]models.StreamConfig, error) {
	streams := map[string]models.StreamConfig{}

	for pair := range txn.IterPairs(kv.Key{"stream", "config"}, kv.KeyRange{}) {
		key := pair.Key()

		stream := key[len(key)-1]

		var streamConfig models.StreamConfig
		if featureflags.GetDemoMode() {
			streamConfig = demoStreamConfig
		} else {
			val := pair.Value()
			// If the value is empty, we treat it as a default config (the zero-value
			// of StreamConfig).
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
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
// "stream:field:<stream>:" existence markers.
func FetchStreamFields(txn kv.QueryTx, stream string) []string {
	fields := []string{}

	for key := range txn.IterKeys(kv.Key{"stream", "field", stream}, kv.KeyRange{}) {
		fieldName := key[len(key)-1]
		fields = append(fields, fieldName)
	}

	return fields
}

// GetOrCreateStreamConfig returns a stream's configuration, lazily creating an
// empty (default) config the first time a stream is referenced so later writes
// have something to update. Demo mode short-circuits to the fixed config.
func GetOrCreateStreamConfig(txn kv.MutationTx, stream string) (models.StreamConfig, error) {
	if featureflags.GetDemoMode() {
		return demoStreamConfig, nil
	}

	var streamConfig models.StreamConfig

	streamKey := kv.Key{"stream", "config", stream}
	val, err := txn.Get(streamKey)
	if err != nil {
		return models.StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)
	}

	if val == nil {
		err := txn.Set(streamKey, []byte{})
		if err != nil {
			return models.StreamConfig{}, fmt.Errorf(
				"could not create default stream config '%s': %w",
				stream, err,
			)
		}
	}

	// If the value is empty, we treat it as a default config (the zero-value
	// of StreamConfig).
	if len(val) > 0 {
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

// SaveStreamConfig stores a new stream configuration and computes how its
// indexed fields differ from the previous one. It returns the fields newly
// added to the index (whose existing records must be back-filled) and the
// fields removed from it (whose index must be dropped).
//
// The actual back-fill and index removal are performed separately, in their own
// batched transactions, because on a large stream they can exceed a single
// transaction's size and time limits.
func SaveStreamConfig(
	txn kv.MutationTx,
	stream string,
	config models.StreamConfig,
) (toIndex []string, toUnindex []string, err error) {
	if config.IndexedFields == nil {
		config.IndexedFields = []string{}
	}

	oldConfig, err := GetOrCreateStreamConfig(txn, stream)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch old stream config '%s': %w", stream, err)
	}

	streamKey := kv.Key{"stream", "config", stream}
	configVal, err := json.Marshal(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
	}

	if err := txn.Set(streamKey, configVal); err != nil {
		return nil, nil, fmt.Errorf("could not save stream config '%s': %w", stream, err)
	}

	for _, field := range config.IndexedFields {
		if !oldConfig.IsFieldIndexed(field) {
			toIndex = append(toIndex, field)
		}
	}

	for _, field := range oldConfig.IndexedFields {
		if !config.IsFieldIndexed(field) {
			toUnindex = append(toUnindex, field)
		}
	}

	return toIndex, toUnindex, nil
}

// DeleteStreamDataBatch deletes up to limit of a stream's data keys — its log
// entries ("entry:<stream>:*"), inverted indexes ("index:<stream>:*") and field
// markers ("stream:field:<stream>:*") — and reports how many it deleted. A count
// of zero means all data is gone (the configuration key is removed separately by
// [DeleteStreamConfig]). Callers loop until it returns zero so the work stays
// within the backend's transaction limits.
func DeleteStreamDataBatch(txn kv.MutationTx, stream string, limit int) (int, error) {
	prefixes := []kv.Key{
		{"entry", stream},
		{"index", stream},
		{"stream", "field", stream},
	}

	var keys []kv.Key
	for _, prefix := range prefixes {
		for key := range txn.IterKeys(prefix, kv.KeyRange{}) {
			keys = append(keys, key)
			if len(keys) >= limit {
				break
			}
		}
		if len(keys) >= limit {
			break
		}
	}

	for _, key := range keys {
		if err := txn.Clear(key); err != nil {
			return len(keys), fmt.Errorf(
				"could not delete key '%s' from stream '%s': %w",
				key, stream, err,
			)
		}
	}

	return len(keys), nil
}

// DeleteStreamConfig removes a stream's configuration key. It is cleared last,
// after all the stream's data, so an interrupted deletion can be resumed.
func DeleteStreamConfig(txn kv.MutationTx, stream string) error {
	streamKey := kv.Key{"stream", "config", stream}
	if err := txn.Clear(streamKey); err != nil {
		return fmt.Errorf("could not delete stream config '%s': %w", stream, err)
	}

	return nil
}
